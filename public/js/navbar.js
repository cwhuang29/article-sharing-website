const logoutEndpoint = "/logout";

const getCookie = (name) => {
  var dc = document.cookie;
  var prefix = name + "=";
  var begin = dc.indexOf("; " + prefix);
  if (begin == -1) {
    begin = dc.indexOf(prefix);
    if (begin != 0) return null;
  } else {
    begin += 2;
    var end = document.cookie.indexOf(";", begin);
    if (end == -1) {
      end = dc.length;
    }
  }
  return decodeURI(dc.substring(begin + prefix.length, end));
};

const logout = () => {
  fetch(logoutEndpoint, {
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    credentials: "same-origin",
    redirect: "follow",
    referrerPolicy: "no-referrer",
  }).then((resp) => {
    if (resp.status >= 400) {
      errMsg = "<div><p><strong>Some Severe Errors Occurred</strong></p><p>Please reload the page and try again.</p></div>";
      showErrMsg(errMsg);
    } else {
      document.cookie = "login_email=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"; // Set the expires parameter to a passed date to delete a cookie
      document.cookie = "login_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      document.cookie = "is_admin=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      [...resp.headers.entries()].forEach((header) => console.log(header[0], header[1]));
      window.location.href = resp.headers.get("Location");
    }
  });
};

const showNewPostButton = () => {
  if (getCookie("is_admin")) {
    newPostBtn.style.display = "block";
  }
};

document.addEventListener("DOMContentLoaded", () => {
  const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll(".navbar-burger"), 0);

  if ($navbarBurgers.length > 0) {
    $navbarBurgers.forEach((el) => {
      el.addEventListener("click", () => {
        const target = el.dataset.target;
        const $target = document.getElementById(target);
        el.classList.toggle("is-active");
        $target.classList.toggle("is-active");
      });
    });
  }

  loginBtn = document.getElementById("loginBtn");
  logoutBtn = document.getElementById("logoutBtn");
  newPostBtn = document.getElementById("newPostBtn");

  if (getCookie("login_email")) {
    loginBtn.style.display = "none";
  } else {
    logoutBtn.style.display = "none";
  }
  logoutBtn.addEventListener("click", logout);
  showNewPostButton();
});
