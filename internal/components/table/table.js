(function () {
  "use strict";

  document.addEventListener("click", function (e) {
    const btn = e.target.closest("span[data-column-sortable='true']");
    if (!btn) return;

    const table = btn.closest("table");
    if (!table) return;

    const tbody = table.querySelector("tbody");
    if (!tbody) return;

    const column = btn.dataset.columnName;
    if (!column) return;

    const currentOrder = btn.dataset.columnOrder === "asc" ? "asc" : "desc";
    const newOrder = currentOrder === "asc" ? "desc" : "asc";

    btn.dataset.columnOrder = newOrder;

    // reset other headers in this table
    table.querySelectorAll("span[data-column-sortable='true']").forEach((b) => {
      if (b !== btn) {
        b.dataset.columnOrder = "";
      }
    });

    const rows = Array.from(tbody.querySelectorAll("tr"));

    rows.sort((a, b) => {
      const aCell = a.querySelector(`td[data-column="${column}"]`);
      const bCell = b.querySelector(`td[data-column="${column}"]`);

      const aValue = aCell?.innerText.trim() ?? "";
      const bValue = bCell?.innerText.trim() ?? "";

      // naive... can only handle strings or numeric strings
      // better approach is needed for values like
      // "1,000", "50%", "2026-3-3", "$93.12"
      if (!isNaN(aValue) && !isNaN(bValue)) {
        return newOrder === "asc"
          ? Number(aValue) - Number(bValue)
          : Number(bValue) - Number(aValue);
      }

      return newOrder === "asc"
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue);
    });

    rows.forEach((row) => tbody.appendChild(row));
  });
})();
