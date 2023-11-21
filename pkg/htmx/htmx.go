package htmx

import (
	"net/http"
	"strconv"
	"strings"
)

func Join(s ...string) string {
	return strings.Join(s, ", ")
}

// GetBoosted indicates that the request is via an element using hx-boost.
func GetBoosted(r *http.Request) bool {
	truthy, _ := strconv.ParseBool(r.Header.Get("HX-Boosted"))
	return truthy
}

// GetCurrentURL is the current URL of the browser.
func GetCurrentURL(r *http.Request) string {
	return r.Header.Get("HX-Current-URL")
}

// GetHistoryRestoreRequest returns true if the request is for history restoration after a miss in the local history cache.
func GetHistoryRestoreRequest(r *http.Request) bool {
	truthy, _ := strconv.ParseBool(r.Header.Get("HX-History-Restore-Request"))
	return truthy
}

// GetPrompt is the user response to an hx-prompt.
func GetPrompt(r *http.Request) string {
	return r.Header.Get("HX-Prompt")
}

// GetRequest indicates that it is a HTMX request.
func GetRequest(r *http.Request) bool {
	truthy, _ := strconv.ParseBool(r.Header.Get("HX-Request"))
	return truthy
}

// GetTarget is the id of the target element if it exists.
func GetTarget(r *http.Request) string {
	return r.Header.Get("HX-Target")
}

// GetTriggerName is the name of the triggered element if it exists.
func GetTriggerName(r *http.Request) string {
	return r.Header.Get("HX-Trigger-Name")
}

// GetTrigger is the id of the triggered element if it exists.
func GetTrigger(r *http.Request) string {
	return r.Header.Get("HX-Trigger")
}

// SetLocation allows you to do a client-side redirect that does not do a full page reload.
func SetLocation(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Location", url)
}

// SetPushURL pushes a new url into the history stack.
func SetPushURL(w http.ResponseWriter) {
	w.Header().Set("HX-Push-URL", "true")
}

// SetRedirect used to do a client-side redirect to a new location.
func SetRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

// SetRefresh will do a client-side full refresh of the page.
func SetRefresh(w http.ResponseWriter) {
	w.Header().Set("HX-Refresh", "true")
}

// SetReplaceURL replaces the current URL in the location bar.
func SetReplaceURL(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Replace-Url", url)
}

// SetReswap allows you to specify how the response will be swapped.
func SetReswap(w http.ResponseWriter) {
	w.Header().Set("HX-Reswap", "true")
}

// SetRetarget is a CSS selector that updates the target of the content update to a different element on the page.
func SetRetarget(w http.ResponseWriter, target string) {
	w.Header().Set("HX-Retarget", target)
}

// SetReselect is a CSS selector that allows you to choose which part of the response is used to be swapped in.
func SetReselect(w http.ResponseWriter) {
	w.Header().Set("HX-Reselect", "true")
}

// SetTrigger events as soon as the response is received.
func SetTrigger(w http.ResponseWriter, event string) {
	w.Header().Set("HX-Trigger", event)
}

// SetTriggerAfterSettle events after the settling step.
func SetTriggerAfterSettle(w http.ResponseWriter, event string) {
	w.Header().Set("HX-Trigger-After-Settle", event)
}

// SetTriggerAfterSwap events after the swap step.
func SetTriggerAfterSwap(w http.ResponseWriter, event string) {
	w.Header().Set("HX-Trigger-After-Swap", event)
}
