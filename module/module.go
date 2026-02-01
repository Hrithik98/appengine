package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/module"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Create a standard App Engine context
	ctx := appengine.NewContext(r)

	// 1. List all modules (services) in the application
	modules, err := module.List(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing modules: %v", err), 500)
		return
	}
	fmt.Fprintf(w, "Modules:\n") // Removed %s as no argument was provided
	for _, m := range modules {
		fmt.Fprintf(w, "- %s\n", m)
	}

	// 2. Identify the current module and its default version
	currentMod := appengine.ModuleName(ctx)
	
	defaultVer, err := module.DefaultVersion(ctx, currentMod)
	if err != nil {
		defaultVer = "unknown"
	}
	fmt.Fprintf(w, "DefaultVersion %s:\n", defaultVer)
	
	// FIX 1: Use module.Versions (plural)
	// FIX 2: Use := to declare the new variable 'versions'
	// FIX 3: Define the module name you want to query
	queryModule := "default"
	versions, err := module.Versions(ctx, queryModule)
	if err != nil {
		http.Error(w, "Could not list versions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// FIX 4: Use 'queryModule' instead of the undefined 'moduleName'
	fmt.Fprintf(w, "Versions for module %s:\n", queryModule)
	for _, v := range versions {
		fmt.Fprintf(w, "- %s\n", v)
	}

	// 3. Scale a specific version (if it uses manual scaling)
	// Note: Removed the trailing space from the module name string
	err = module.SetNumInstances(ctx, "bundled-services-mail-java-sdk-app", "20251016t135918", 5)
	if err != nil {
		fmt.Printf("Scaling failed (expected if not manual scaling): %v\n", err)
	}
}


func main() {
	http.HandleFunc("/", handler)
	appengine.Main()
}

