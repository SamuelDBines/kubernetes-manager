(() => {
  const input = document.getElementById("ns-search");
  const grid = document.getElementById("ns-grid");
  if (!input || !grid) return;

  const tiles = Array.from(grid.querySelectorAll(".tile"));

  input.addEventListener("input", () => {
    const q = input.value.trim().toLowerCase();
    for (const t of tiles) {
      const title = t.querySelector(".tile__title")?.textContent?.toLowerCase() || "";
      t.style.display = title.includes(q) ? "" : "none";
    }
  });
})();
