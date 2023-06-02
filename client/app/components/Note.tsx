import { Note } from "@/types";
import React from "react";

interface Props {
  note: Note;
}

function Note({ note }: Props) {
  return (
    <div>
      <h3>{note.details}</h3>
      <p>{note.timestamp}</p>
    </div>
  );
}

export default Note;
