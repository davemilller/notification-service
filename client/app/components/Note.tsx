import { Note as Notification } from "@/types";
import React from "react";

interface Props {
  note: Notification;
}

function Notification({ note }: Props) {
  return (
    <div>
      <h3>{note.details}</h3>
      <p>{note.timestamp}</p>
    </div>
  );
}

export default Notification;
