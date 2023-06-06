# notification-service

An example service for pushing notifications from server to client using golang and graphql subscriptions.

Driven by clean architecture & best practices.

### Archtecture
* `app/`
  * This folder holds the golang backend code for the graphql server.
  * It is split into a control layer (endpoints), a domain layer (models/structs), and a framework layer with a redis package and a subscription manager package

* `client/`
  * This folder holds the frontend react / nextJS code.

* `config/`
  * This folder has the docker files and the `.env`

### Usage
* From the root, run:
  * `go mod vendor`
  * `make local`
* From `./client`, run:
  * `npm install`
  * `npm run dev`

### Design
The goal for this app is to have a nice example of a real-time server to client notification system using graphql subscriptions, specifically using [graphql-go](https://github.com/graphql-go/graphql) and [apollo](https://www.apollographql.com/docs/react/data/subscriptions/) subscriptions for the frontend. There seems to be a lack of resources on this topic using these packages, especially any self-contained examples that use both server and client.

Thanks to [this](https://github.com/eientei/wsgraphql) websocket implementation, only a small wrapper around the websocket connection is necessary.

The idea is that a client can load the frontend, make a websocket connection (using a unique identifier), and get pushed real time messages.

1. The client queries for outstanding messages.
```graphql
query GetNotes($userID: String!) {
  getNotes(userID: $userID) {
    id
    userID
    details
    timestamp
  }
}
```

2. The client initializes a subscription connection using this graphql and apollo's `subscribeToMore`.
```graphql
subscription Notifications($userID: String!) {
  notifications(userID: $userID) {
    id
    userID
    details
    timestamp
  }
}
```
```js
subscribeToMore({
  document: NOTIFICATIONS_SUBSCRIPTION,
  variables: { userID: userID },
  updateQuery: (prev, { subscriptionData }) => {
    if (!subscriptionData.data) return prev;
    const newNote = subscriptionData.data.notifications;
    setNotes((notes) => [newNote, ...notes]);
  },
});
```

3. Notifications are added to a user's queue using this graphql request. If there is a live connection for this user, the message is pushed to the websocket, otherwise it is cached in a redis db and all outstanding messages are pushed when the user connects.
```graphql
mutation AddNote($userID: String, $details: String) {
    addNote(userID: $userID, details: $details) {
        id
        userID
        details
        timestamp
    }
}
```

4. Messages are ack'd out of the user's queue as they are delivered. An alternative to this would be to save them to a database for audit purposes or have the user manually ack them with a button.

### Future work
* Make a nice UI
* Handle the user ID via context
* Authentication
* Package / microservice?
* Write websocket implementation of my own
