const API = window.location.hostname === "localhost" 
    ? "http://localhost:8080" 
    : "https://algoroulette.up.railway.app";

function getHeaders() {
    return {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + localStorage.getItem("token"),
    };
}

export async function loginOrRegister(mode, email, password) {
    const res = await fetch(`${API}/auth/${mode}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    return { res, data };
}

export async function fetchQuestion(topic, diff) {
    let url = `${API}/question`;
    const params = [];
    if (topic) params.push("topic=" + topic);
    if (diff) params.push("difficulty=" + diff);
    if (params.length) url += "?" + params.join("&");

    const res = await fetch(url, { headers: getHeaders() });
    return res;
}

export async function updateQuestionStatus(id, status) {
    return fetch(`${API}/question/${id}/status`, {
        method: "PATCH",
        headers: getHeaders(),
        body: JSON.stringify({ status }),
    });
}

export async function updateQuestionBookmark(id, bookmarked) {
    return fetch(`${API}/question/${id}/bookmark`, {
        method: "PATCH",
        headers: getHeaders(),
        body: JSON.stringify({ bookmarked }),
    });
}

export async function updateQuestionNotes(id, notes) {
    return fetch(`${API}/question/${id}/notes`, {
        method: "PATCH",
        headers: getHeaders(),
        body: JSON.stringify({ notes }),
    });
}

export async function postQuestion(payload) {
    return fetch(`${API}/question`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(payload),
    });
}

export async function postImport(questions) {
    return fetch(`${API}/questions/import`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(questions),
    });
}

export async function deleteQuestions() {
    return fetch(`${API}/questions`, {
        method: "DELETE",
        headers: getHeaders(),
    });
}