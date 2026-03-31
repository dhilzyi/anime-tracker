// web/app.js

const listEl = document.getElementById("anime-list");

async function fetchAnime() {
  const res = await fetch("/api/anime");
  const data = await res.json();

  listEl.innerHTML = "";

  data.forEach((anime) => {
    const li = document.createElement("li");
    li.textContent = anime.romaji_name;

    const btn = document.createElement("button");
    btn.textContent = "Delete";
    btn.onclick = async () => {
      await fetch(`/api/anime/${anime.id}`, { method: "DELETE" });
      fetchAnime();
    };

    li.appendChild(btn);
    listEl.appendChild(li);
  });
}

async function fetchAnimeData() {
  const res = await fetch("/api/anime");
  const data = await res.json();
  return data;
}
const mainTable = document.getElementById("main-table");

let numberRow = 1 + 1;
async function addDataToTable() {
  const animeData = await fetchAnimeData();
  animeData.forEach((anime) => {
    const tableRow = document.createElement("tr");

    const tdNumber = document.createElement("td");
    tdNumber.textContent = numberRow;
    tableRow.appendChild(tdNumber);

    const columns = [
      "id",
      "romaji_name",
      "japanese_name",
      "english_name",
      "type",
      "release_date",
    ];

    columns.forEach((key) => {
      const td = document.createElement("td");
      td.textContent = anime[key] ?? "";
      tableRow.appendChild(td);
    });
    mainTable.appendChild(tableRow);
    numberRow += 1;
  });
}

async function addOneRowToTable(data) {
  const tableRow = document.createElement("tr");

  const tdNumber = document.createElement("td");
  tdNumber.textContent = numberRow;
  tableRow.appendChild(tdNumber);

  Object.entries(data).map((entry) => {
    const tableData = document.createElement("td");
    let value = entry[1];
    tableData.textContent = value;
    tableRow.appendChild(tableData);
  });
  mainTable.appendChild(tableRow);
  numberRow += 1;
}

document.getElementById("anime-form").onsubmit = async (e) => {
  e.preventDefault();
  const form = e.target;

  const body = {
    romaji_name: form.romaji_name.value,
    japanese_name: form.japanese_name.value || null,
    english_name: form.english_name.value || null,
    type: form.type.value || null,
    release_date: form.release_date.value || null,
  };

  const res = await fetch("/api/anime", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  if (res.ok) {
    form.reset();
    const data = await res.json();
    addOneRowToTable(data);
  }
};

addDataToTable();
