"use strict";

const form = document.getElementById("provider-search");
if (form) {
  const input = document.getElementById("search-query");
  const status = document.getElementById("search-status");
  const words = (text) => text.normalize("NFKC").toLowerCase().match(/[\p{L}\p{N}]+/gu) ?? [];
  const rows = Array.from(document.querySelectorAll("#provider-list > li"), (node) => ({
    node,
    terms: new Set(words(node.dataset.searchText)),
  }));

  function showResults(query) {
    const terms = words(query);
    let count = 0;
    for (const {node, terms: profile} of rows) {
      const matches = terms.every((term) => profile.has(term));
      node.hidden = !matches;
      if (matches) count++;
    }
    status.textContent = count === 0 ? "No matching providers" : `${count} ${count === 1 ? "provider" : "providers"}`;
  }

  input.addEventListener("input", () => showResults(input.value));
  form.addEventListener("submit", (event) => event.preventDefault());
  form.addEventListener("reset", (event) => {
    event.preventDefault();
    input.value = "";
    input.focus();
    showResults("");
  });
  form.hidden = false;
}
