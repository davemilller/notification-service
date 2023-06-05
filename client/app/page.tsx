"use client";

import { useEffect, useState } from "react";
import { useSubscription } from "@apollo/client";
import { NOTIFICATIONS_SUBSCRIPTION } from "@/api";
import Note from "./components/Note";

export default function Home() {
  const [userID, setUserID] = useState("0");
  const [notes, setNotes] = useState<Note[]>([]);

  const { data, error } = useSubscription(NOTIFICATIONS_SUBSCRIPTION, {
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
      setNotes((notes) => [...notes, data.notifications]);
    }
  }, [data, userID]);

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        padding: "10px",
      }}
    >
      <h1>Notes for user: {userID}</h1>
      <div style={{ display: "flex", flexDirection: "column" }}>
        {notes.map((note) => (
          <Note key={note.id} note={note} />
        ))}
      </div>
    </div>
  );
}
