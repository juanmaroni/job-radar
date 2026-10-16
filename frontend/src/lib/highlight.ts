import { buildRegEx } from "@/lib/utils";

export function highlightHtml(html: string, query: string): string {
  if (query === "") return html;

  const divContainer = document.createElement("div");
  divContainer.innerHTML = html;

  const nodes = document.createTreeWalker(divContainer, NodeFilter.SHOW_TEXT);
  const nodesArr: Text[] = [];
  let node: Text | null;

  while ((node = nodes.nextNode() as Text | null)) {
    nodesArr.push(node);
  }

  const regex: RegExp = buildRegEx(query);
  
  for (const n of nodesArr) {
    const value = n.nodeValue;

    if (value === null) continue;

    const groups = value.split(regex);

    if (groups.length < 2) continue;

    const doc = document.createDocumentFragment();

    for (const [index, text] of groups.entries()) {
      if (index % 2 === 0) {
        doc.appendChild(document.createTextNode(text));
      } else {
        const markElem = document.createElement("mark");
        const textElem = document.createTextNode(text);
        markElem.appendChild(textElem);
        doc.appendChild(markElem);
      }
    }

    n.replaceWith(doc);
  }

  return divContainer.innerHTML;
}
