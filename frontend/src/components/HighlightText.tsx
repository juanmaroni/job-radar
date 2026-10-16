"use client";

import { buildRegEx } from "@/lib/utils";

export function HighlightText({text, query}: {text: string; query: string}) {
  if (query === "") return text;

  const groups = text.split(buildRegEx(query));

  return (
    <>
      {groups.map((e, index) =>
        index % 2 === 1 ? (
          <mark key={index}>{e}</mark>
        ) : (
          <span key={index}>{e}</span>
        )
      )}
    </>
  );
}
