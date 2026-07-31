# Anthologise

## Dev Info

### Running the Project

Taskfile is used to manage developement commands and different lifecycle targets.
The root taskfile acts as an entrypoint to sub-taskfiles. Each sub-taskfile manages the compose commands for its container (and uses the workdir of the super taskfile).

To see available tasks:
```bash
task
```

To see available tasks for a give sub-taskfile, for example service:
```bash
task service
```

Once the service is running, to add the locally hosted manifest to Stremio, paste the following into Stremio's add by URL option:
```text
http://127.0.0.1:7000/manifest.json
```

### Env Files

All dev-facing configuration is handled through environment variables. These are pulled in through the Viper library in the service's config package.

In the service directory, `.env.example` shows an example of all the variables expected to be set for the project. 

To set the local config run:
```bash
task service:envinit
```
Then modify any variables in `.env.local`.