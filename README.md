```
 _   _      _                _                 _   
 | | | | ___| |_ __ ___      / \   ___ ___  ___| |_ 
 | |_| |/ _ \ | '_ ` _ \    / _ \ / __/ __|/ _ \ __|
 |  _  |  __/ | | | | | |  / ___ \\__ \__ \  __/ |_ 
 |_| |_|\___|_|_| |_| |_| /_/   \_\___/___/\___|\__|
```

# helm-plugin-asset
A helm plugin to manage chart assets

Usage
=====

This is meant to be invoked by [helm](https://github.com/kubernetes/helm) as helm plugin. To install:

    helm plugin install ... # TODO


Then run:

    helm asset --asset-dir assets --values values.yaml

This will recursively render the asset files under `assets` folder, and populate the `data` key of the values file with the rendered content.
Then you can install a chart with the new `values.yaml` file.
