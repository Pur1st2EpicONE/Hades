const api = "/api/v1";

async function loadItems(query = "") {
  const res = await fetch(`${api}/items${query}`);
  const json = await res.json();
  const items = json.result || [];
  const tbody = document.querySelector("#itemsTable tbody");
  tbody.innerHTML = "";
  items.forEach((item) => {
    const tr = document.createElement("tr");
    const isIncome = item.type === "income";
    tr.innerHTML = `
            <td>${item.id}</td>
            <td><span class="${isIncome ? "income" : "expense"}">${item.type}</span></td>
            <td>${parseFloat(item.amount).toFixed(2)}</td>
            <td>${new Date(item.date).toLocaleDateString("ru-RU")}</td>
            <td>${item.category}</td>
            <td>${item.description || ""}</td>
            <td><button data-id="${item.id}" class="deleteBtn">🗑</button></td>
        `;
    tbody.appendChild(tr);
  });
  bindDeleteButtons();
}

function bindDeleteButtons() {
  document.querySelectorAll(".deleteBtn").forEach((btn) => {
    btn.onclick = async () => {
      await fetch(`${api}/items/${btn.dataset.id}`, { method: "DELETE" });
      loadItems();
    };
  });
}

document.getElementById("createForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const payload = Object.fromEntries(formData);
  const res = await fetch(`${api}/items`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (res.ok) {
    e.target.reset();
    loadItems();
  }
});

document.getElementById("filterForm").addEventListener("submit", (e) => {
  e.preventDefault();
  const params = new URLSearchParams(new FormData(e.target));
  loadItems("?" + params.toString());
});

document.getElementById("exportCSV").onclick = () => {
  const params = new URLSearchParams(
    new FormData(document.getElementById("filterForm")),
  );
  params.set("export", "csv");
  window.location = `${api}/items?${params}`;
};

document
  .getElementById("analyticsForm")
  .addEventListener("submit", async (e) => {
    e.preventDefault();
    const params = new URLSearchParams(new FormData(e.target));
    const res = await fetch(`${api}/analytics?${params}`);
    const json = await res.json();
    const data = json.result;
    const container = document.getElementById("analyticsResult");
    if (Array.isArray(data)) {
      let html = `<table class="grouped"><thead><tr>
        <th>Group</th>
        <th>Count</th>
        <th>Income</th>
        <th>Expense</th>
        <th>Balance</th>
        <th>Avg Amount</th>
      </tr></thead><tbody>`;
      data.forEach((row) => {
        const balanceClass =
          parseFloat(row.balance) >= 0 ? "income" : "expense";
        html += `
                <tr>
                    <td>${row.group_key}</td>
                    <td>${row.count}</td>
                    <td>${parseFloat(row.total_income).toFixed(2)}</td>
                    <td>${parseFloat(row.total_expense).toFixed(2)}</td>
                    <td class="${balanceClass}">${parseFloat(row.balance).toFixed(2)}</td>
                    <td>${parseFloat(row.avg_amount).toFixed(2)}</td>
                </tr>`;
      });
      html += `</tbody></table>`;
      container.innerHTML = html;
    } else {
      const balanceClass = parseFloat(data.balance) >= 0 ? "income" : "expense";
      container.innerHTML = `
            <div class="stats">
                <div>Records: <strong>${data.count}</strong></div>
                <div>Income: <strong class="income">${parseFloat(data.total_income).toFixed(2)}</strong></div>
                <div>Expense: <strong class="expense">${parseFloat(data.total_expense).toFixed(2)}</strong></div>
                <div>Balance: <strong class="${balanceClass}">${parseFloat(data.balance).toFixed(2)}</strong></div>
                <div>Average: ${parseFloat(data.avg_amount).toFixed(2)}</div>
                <div>Median: ${parseFloat(data.median).toFixed(2)}</div>
                <div>P90: ${parseFloat(data.percentile_90).toFixed(2)}</div>
            </div>`;
    }
  });

document.getElementById("exportAnalytics").onclick = () => {
  const params = new URLSearchParams(
    new FormData(document.getElementById("analyticsForm")),
  );
  params.set("export", "csv");
  window.location = `${api}/analytics?${params}`;
};

loadItems();
