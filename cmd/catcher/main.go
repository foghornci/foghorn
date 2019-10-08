package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	apiv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn.jenkins.io/v1alpha1"
	clientv1alpha1 "github.com/foghornci/foghorn/pkg/client/clientset/versioned/typed/foghorn.jenkins.io/v1alpha1"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/bitbucket"
	"github.com/jenkins-x/go-scm/scm/driver/gitea"
	"github.com/jenkins-x/go-scm/scm/driver/github"
	"github.com/jenkins-x/go-scm/scm/driver/gitlab"
	"github.com/jenkins-x/go-scm/scm/driver/gogs"
	"github.com/jenkins-x/go-scm/scm/driver/stash"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

const (
	helloMessage = "hello from Foghorn!\n"

	// HealthPath is the URL path for the HTTP endpoint that returns health status.
	HealthPath = "/health"
	// ReadyPath URL path for the HTTP endpoint that returns ready status.
	ReadyPath = "/ready"
)

// WebhookOptions holds the command line arguments
type WebhookOptions struct {
	BindAddress   string
	Path          string
	Port          int
	JSONLog       bool
	GitClient     *scm.Client
	FoghornClient *clientv1alpha1.FoghornV1alpha1Client
	KubeConfig    *rest.Config

	namespace string
}

func main() {

	gitClient, err := getGitClient()
	if err != nil {
		logrus.WithError(err).Fatal("could not initialize git client: %s", err)
	}

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		logrus.WithError(err).Fatal("could not initialize k8s client: %s", err)
	}

	foghornClient, err := clientv1alpha1.NewForConfig(kubeConfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not initialize Foghorn client: %s", err)
	}

	o := WebhookOptions{
		Path:          "/",
		Port:          8080,
		JSONLog:       true,
		BindAddress:   "localhost",
		GitClient:     gitClient,
		KubeConfig:    kubeConfig,
		FoghornClient: foghornClient,
	}

	if o.JSONLog {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	mux := http.NewServeMux()
	mux.Handle(HealthPath, http.HandlerFunc(o.health))
	mux.Handle(ReadyPath, http.HandlerFunc(o.ready))

	mux.Handle(o.Path, http.HandlerFunc(o.handleWebHookRequests))

	logrus.Infof("catcher is now listening on path %s for webhooks from %s", o.Path, o.GitClient.Driver.String())
	if err := http.ListenAndServe(":"+strconv.Itoa(o.Port), mux); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// health returns either HTTP 204 if the service is healthy, otherwise nothing ('cos it's dead).
func (o *WebhookOptions) health(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Health check")
	w.WriteHeader(http.StatusNoContent)
}

// ready returns either HTTP 204 if the service is ready to serve requests, otherwise HTTP 503.
func (o *WebhookOptions) ready(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Ready check")
	if o.isReady() {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// getIndex returns a simple home page
func (o *WebhookOptions) getIndex(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("GET index")
	w.Write([]byte(helloMessage))
}

func (o *WebhookOptions) isReady() bool {
	// TODO a better readiness check
	return true
}

// handle request for pipeline runs
func (o *WebhookOptions) handleWebHookRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		o.handleGet(w, r)
	case http.MethodPost:
		o.handlePost(w, r)
	default:
		logrus.Infof("catcher received and ignored a request with unsupported method %s from %s", r.Method, r.RemoteAddr)
	}
}

func getGitClient() (*scm.Client, error) {

	gitKind := strings.ToLower(os.Getenv("GIT_KIND"))
	gitURL := os.Getenv("GIT_URL")

	switch gitKind {
	case "bitbucket", "bitbucket-cloud", "bitbucketcloud":
		if gitURL != "" {
			return bitbucket.New(gitURL)
		}
		return bitbucket.NewDefault(), nil
	case "bitbucketserver", "bitbucket-server":
		return stash.New(gitURL)
	case "gitea":
		return gitea.New(gitURL)
	case "gitlab":
		return gitlab.New(gitURL)
	case "gogs":
		return gogs.New(gitURL)
	default:
		if gitURL != "" {
			return github.New(gitURL)
		}
		return github.NewDefault(), nil
	}
}

func secretFunc(scm.Webhook) (string, error) {
	token := os.Getenv("HMAC_TOKEN")
	if token == "" {
		return "", fmt.Errorf("hmac token not found")
	}
	return token, nil
}

func (o *WebhookOptions) handleGet(w http.ResponseWriter, r *http.Request) {
}

func (o *WebhookOptions) handleAuthenticatedPost(w http.ResponseWriter, r *http.Request, hook scm.Webhook) {

	// Initialize GitEvent
	gitEvent := &apiv1alpha1.GitEvent{
		Spec: apiv1alpha1.GitEventSpec{
			EventType: "generic",
			ParsedWebhook: scm.WebhookSerializer{
				Webhook: hook,
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "gitevent-",
		},
	}

	gitEventInterface := o.FoghornClient.GitEvents("foghorn")

	// Create GitEvent resource
	result, err := gitEventInterface.Create(gitEvent)
	if err != nil {
		logrus.WithError(err).Fatal("GitEvent CRD creation failed: %s", err)
	}
	repo := result.Spec.ParsedWebhook.Webhook.Repository()
	logrus.Infof("GitEvent CRD created for repo %s/%s", repo.Name, repo.Namespace)
	w.Write([]byte("Thanks for the webhook!"))
}

func (o *WebhookOptions) handlePost(w http.ResponseWriter, r *http.Request) {
	// Parse webhook
	parsedWebhook, err := o.GitClient.Webhooks.Parse(r, secretFunc)
	if err != nil {
		logrus.Warnf("error during webhook parsing: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		o.handleAuthenticatedPost(w, r, parsedWebhook)
	}

}

func (o *WebhookOptions) returnError(err error, message string, w http.ResponseWriter, r *http.Request) {
	logrus.Errorf("returning error: %v %s", err, message)
	responseHTTPError(w, http.StatusInternalServerError, "500 Internal Error: "+message+" "+err.Error())
}

func responseHTTPError(w http.ResponseWriter, statusCode int, response string) {
	logrus.WithFields(logrus.Fields{
		"response":    response,
		"status-code": statusCode,
	}).Info(response)
	http.Error(w, response, statusCode)
}
