"use client";

import { useEffect, useState } from "react";
import { useQuery, useSubscription } from "@apollo/client";
import { GET_NOTIFICATIONS, NOTIFICATIONS_SUBSCRIPTION } from "@/api";
import Notification from "./components/Note";
import Notifications from "./components/Notifications";

export default function Home() {
  const [userID, setUserID] = useState("0");
  const [notes, setNotes] = useState<Notification[]>([]);

  const { subscribeToMore, data, error } = useQuery(GET_NOTIFICATIONS, {
    variables: {
      userID: userID,
    },
  });
  if (error) {
    console.error(error);
  }

  console.log(data);

  useEffect(() => {
    if (data) {
      setNotes(data.getNotes);
    }
  }, [data, userID]);

  return (
    <Notifications
      userID={userID}
      notes={...notes}
      subscribe={() => {
        subscribeToMore({
          document: NOTIFICATIONS_SUBSCRIPTION,
          variables: { userID: userID },
          updateQuery: (prev, { subscriptionData }) => {
            if (!subscriptionData.data) return prev;
            const newNote = subscriptionData.data.notifications;
            setNotes((notes) => [newNote, ...notes]);
          },
        });
      }}
    />
  );
}
