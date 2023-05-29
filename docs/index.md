# Welcome to KIE (Knowledge Is Everything)!

A collection of microservices to predict who will be good at hockey.

## Background

Looking for talk ideas for gophercon, I asked chatGPT what talk ideas might be good.
It listed 10.
I chose recommendation services, as it's a thing I've tried in the past, but never successfully implemented.
After some refining, I decided to make a recommendation service for hockey players.
ChatGPT helped breakdown the problem into 4 major steps, check [outline](pipeline-outline.md) for more info.

## Tools

- Go, absolutely everywhere - because why not?
- Pachyderm - because it's awesome
- gRPC/buf

## Deploying

### Local

```sh
# Start pachyderm
helm install pachd pachyderm/pachyderm -n pachd --create-namespace \
  --set deployTarget=LOCAL \
  --set proxy.enabled=false

# Run the pachyderm operator (github.com/alexanderjophus/pachyderm-operator)

# Apply kustomize files
kubectl apply -k deploy/base
```

### Prod

lol