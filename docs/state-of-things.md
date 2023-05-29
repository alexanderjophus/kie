# State of Things

Using the model described at the bottom of the page to describe how 'done' a component is.

## Components

### UI Input

5 - I'd like a CLI tool for this at first. Potentially uploading data player by player, as well as a 'star' flag for each player.

### Ingestion

50 - We have a gRPC server that can receive data, it then puts that data into a pachyderm repo

### Cleaning

50 - We have a service that converts a json response into a tidy csv with relevant data only. It needs tests, and some tweaks to the data it's outputting.

### Features

10 - We have a service that can take a csv and output a csv with features. It's very hacked together, need to read up on features again too.

### Training

5 - I'd like a job that can be triggered on demand that will train a model. This way we're not training on every commit. I'd like to start with a simple linear regression model. It also needs to store the model somewhere.

### Predictor

5 - After the model has been trained, this service needs to be redeployed to pick up the model. Given a players ID, it should look up that players data so far and predict their future.

## 5/10/50/MVP Model

### 5

The component has been envisioned.

- The goal is somewhat documented.

### 10

The component has been designed.

- The goal is documented.
- The core component has been developed.
- The component has **not** necessarily been tested, neither manually, or through automated tests.

### 50

The component is structured.

- The goal is documented.
- The core component has been developed.
- The component has been tested, either manually, or through automated tests.
- The component has been documented.

### MVP

The component is 'ready' by all accounts.

- The component has been developed.
- The component has been tested through automated tests.
- The component has been documented thoroughly.