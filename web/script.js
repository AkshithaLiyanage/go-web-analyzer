document.getElementById("analyze-form").addEventListener("submit", async function (e) {
  e.preventDefault();

  const urlInput = document.getElementById("url-input").value;
  const resultsDiv = document.getElementById("results");
  const errorDiv = document.getElementById("error");
  const loading = document.getElementById("loading");

  resultsDiv.classList.add("hidden");
  errorDiv.classList.add("hidden");
  loading.classList.remove("hidden");

  const formData = new FormData();
  formData.append("url", urlInput);

  try {
    const response = await fetch("/analyze", {
      method: "POST",
      body: formData,
    });

    const data = await response.json();
    loading.classList.add("hidden");

    if (!response.ok) {
      errorDiv.textContent = `${data.status_code} - ${data.error || "Error analyzing page"}`;
      errorDiv.classList.remove("hidden");
      return;
    }

    let html = `<h2>Analysis Result</h2>`;
    html += `<p><strong>Title:</strong> ${data.title}</p>`;
    html += `<p><strong>HTML Version:</strong> ${data.html_version}</p>`;

    html += `<h3>Headings</h3><ul>`;
    for (const [key, value] of Object.entries(data.headings)) {
      html += `<li>${key.toUpperCase()}: ${value}</li>`;
    }
    html += `</ul>`;

    html += `<p><strong>Login Form Present:</strong> ${data.is_login_form}</p>`;

    html += renderLinksSection("Internal Links", data.internal_links);


    html += renderLinksSection("External Links", data.external_links);

    html += renderLinksSection("Inaccessible Links", data.inaccessible_links);


    resultsDiv.innerHTML = html;
    resultsDiv.classList.remove("hidden");

  } catch (err) {
    loading.classList.add("hidden");
    errorDiv.textContent = `Something went wrong: ${err.message}`;
    errorDiv.classList.remove("hidden");
  }
});

function renderLinksSection(title, links) {
  let html = ""
  if (Array.isArray(links) && links.length > 0) {
     const sectionId = title.replace(/\s+/g, "-").toLowerCase();
    html += `<h3>${title} (${links.length})</h3>`;
    html += `<button onclick="toggleLinks('${sectionId}')">Show Links</button>`;
    html += `<ul id="${sectionId}" style="display: none;">`;
    links.forEach(link => {
      html += `<li><a href="${link}" target="_blank">${link}</a></li>`;
    });
    html += `</ul>`;
  } else {
    html += `<h3>${title} (0)</h3><p>None found.</p>`;
  }

  return html;
}

function toggleLinks(id) {
  const element = document.getElementById(id);
  if (element.style.display === "none") {
    element.style.display = "block";
  } else {
    element.style.display = "none";
  }
}

