const logoutEndpoint = '/logout';

const getCookie = (name) => {
  var dc = document.cookie;
  var prefix = name + '=';
  var begin = dc.indexOf('; ' + prefix);
  if (begin == -1) {
    begin = dc.indexOf(prefix);
    if (begin != 0) return null;
  } else {
    begin += 2;
    var end = document.cookie.indexOf(';', begin);
    if (end == -1) {
      end = dc.length;
    }
  }
  return decodeURIComponent(dc.substring(begin + prefix.length, end));
};

const logout = () => {
  fetchData(logoutEndpoint, {
    method: 'POST',
    cache: 'no-cache',
    credentials: 'same-origin',
    redirect: 'follow',
  }).then((resp) => {
    if (resp.status >= 400) {
      showErrMsg('An Error Occurred', 'Please reload the page and try again.');
    } else {
      // Copy from overview.js
      window.sessionStorage.removeItem('offset');
      window.sessionStorage.removeItem('overviewContent');

      // It is okay to leave `login_token` and `csrf_token` untouched (they've been set httpOnly)
      document.cookie = "login_email=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"; // Set the expires parameter to a passed date to delete a cookie
      document.cookie = "is_admin=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

      window.location.href = resp.headers.get('Location');
    }
  });
};

const showNewPostButton = () => {
  newPostBtn.style.display = getCookie('is_admin') ? 'block' : 'none';
};

const showLoginOrLogout = () => {
  if (getCookie('login_email')) {
    logoutParent.style.display = 'block';
    homeBtn.style.display = 'block';
  } else {
    loginBtn.style.display = 'block';
  }
};

const navbarHandler = () => {
  const navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

  if (navbarBurgers.length > 0) {
    navbarBurgers.forEach((el) => {
      el.addEventListener('click', () => {
        const target = el.dataset.target;
        const $target = document.getElementById(target);
        el.classList.toggle('is-active');
        $target.classList.toggle('is-active');
      });
    });
  }

  loginBtn = document.getElementById('loginBtn');
  logoutParent = document.getElementById('logoutParent');
  logoutSection = document.getElementById('logoutSection');
  homeBtn = document.getElementById('homeBtn');
  newPostBtn = document.getElementById('newPostBtn');

  showLoginOrLogout();
  showNewPostButton();
  logoutSection.addEventListener('click', logout);
};

document.addEventListener('DOMContentLoaded', () => {
  navbarHandler();
});
