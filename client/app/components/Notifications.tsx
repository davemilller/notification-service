import { Note } from "@/types";
import React, { useEffect } from "react";
import Notification from "./Note";

interface Props {
  userID: string;
  notes: Note[];
  subscribe: () => void;
}

function Notifications({ userID, notes, subscribe }: Props) {
  useEffect(() => subscribe(), [subscribe]);

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
          <Notification key={note.id} note={note} />
        ))}
      </div>
    </div>
  );
}

export default Notifications;
