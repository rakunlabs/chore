type tocT = {
  href?: string;
  name: string;
  indent: number;
};

const generateToc = (v: HTMLElement) => {
  const articles: tocT[] = [];

  for (const child of v.children) {
    const a = child.querySelector("a");

    switch (child.nodeName.toLowerCase()) {
    case "h1":
    case "h2":
    case "h3":
      articles.push({
        href: a?.href,
        name: child.textContent,
        indent: ~~child.nodeName[1],
      } as tocT);
      break;

    default:
      break;
    }
  }

  return articles;
};

export { generateToc };
export type { tocT };
