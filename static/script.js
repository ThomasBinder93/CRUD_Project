
    const API = "/api/items";

    async function loadItems() {
        const res = await fetch(API);
        const json = await res.json();
        const items = Array.isArray(json) ? json : json?.data ?? [];

        const list = get("list");
        list.innerHTML = "";

        items.forEach(item => {
            const tr = createElement("tr");

            tr.innerHTML = `<td><input id="${item.id}" value="${item.name}"></td>
                            <td><button onclick="updateItem(${item.id}, getValue(${item.id}))">Ändern</button></td>
                            <td><button onclick="deleteItem(${item.id})">Löschen</button></td>       
            `;
            list.appendChild(tr);
        });
    }

    async function createItem() {

        const name = getValue("nameInput");

        const res = await fetch(API, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name})
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }

        get("nameInput").value = "";
        loadItems();
        showMessage("Item erfolgreich erstellt");
    }

    async function updateItem(id, name) {
        console.log("Id ", id, " Name ", name);

        const res = await fetch(`${API}/${id}`, {
            method: "PUT",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name})
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }


        loadItems();
        showMessage("Item erfolgreich aktualisiert");
    }

    async function deleteItem(id) {
        const res = await fetch(`${API}/${id}`, {
            method: "DELETE"
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }

        loadItems();
        showMessage("Item erfolgreich gelöscht");
    }

    loadItems();

