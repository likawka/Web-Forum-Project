function copyLink(postId) {
  // Construct the URL based on the post ID
  var postUrl = window.location.origin + postId;

  // Copy the URL to the clipboard
  navigator.clipboard.writeText(postUrl)
    .then(function () {
      showNotification("Copied!", "success");
    })
    .catch(function (error) {
      showNotification("Error copying URL: " + error, "error");
    });
}

function showNotification(message, type) {
  var notification = document.getElementById("notification");
  notification.textContent = message;
  notification.className = "show " + type;
  setTimeout(function () {
    notification.className = notification.className.replace("show", "");
  }, 1000);
}


document.getElementById("search").addEventListener("focus", function () {
    var searchHints = document.getElementById("searchHints");
    searchHints.style.display = "block";
});

document.addEventListener("click", function (event) {
    var searchHints = document.getElementById("searchHints");
    if (!event.target.matches('.searchInput')) {
        searchHints.style.display = "none";
    }
});


document.addEventListener('DOMContentLoaded', function() {
  const sidebarToggle = document.querySelector('.sidebar-toggle');
  const sideBar = document.querySelector('.side-bar');

  sidebarToggle.addEventListener('click', function() {
      sideBar.classList.toggle('active');
  });
});