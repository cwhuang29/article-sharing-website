const loginEndpoint = "/login";

const validateEmail = (email) => {
  const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
};

const getLoginFields = async (email, password) => {
  let canSubmit = true;

  if (email.length == 0) {
    canSubmit = false;
    err_msg_email.innerText = "The email can't be empty.";
  } else if (!validateEmail(email)) {
    canSubmit = false;
    err_msg_email.innerText = "Please fill in the correct email.";
  } else {
    err_msg_email.innerText = "";
  }

  if (password.length == 0) {
    canSubmit = false;
    err_msg_password.innerText = "The password can't be empty.";
  } else {
    err_msg_password.innerText = "";
  }

  if (canSubmit) {
    return Promise.resolve();
  }
  return Promise.reject();
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  } else {
    return Promise.resolve(resp);
  }
};

const loginSucceed = async (resp) => {
  window.location.href = resp.headers.get("Location");
  return Promise.resolve();
};

const loginFailed = async (resp) => {
  if (resp.status == 500) {
    resp.json().then(function (data) {
      showErrMsg(`<div><p><strong>Error !</strong></p><p>${data.err}</p></div>`);
    });
  } else {
    resp.json().then(function (data) {
      c(data);
      showErrMsg(data.err);
      for (var key in data.errTags) {
        document.getElementById(`err_msg_${key}`).innerText = data.err[key];
      }
    });
  }
};

const postData = async (url = "", data = {}) => {
  return fetch(url, {
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    credentials: "same-origin",
    headers: {
      "Content-Type": "application/json",
    },
    redirect: "follow",
    referrerPolicy: "no-referrer",
    body: JSON.stringify(data),
  });
};

const login = async () => {
  submitBtn.classList.add("is-loading");

  let email = email_field.value.trim();
  let password = password_field.value.trim();
  let res = await getLoginFields(email, password)
    .then(() => {
      return 1;
    })
    .catch(() => {
      return 0;
    });

  if (!res) {
    submitBtn.classList.remove("is-loading");
    return;
  }
  postData(loginEndpoint, {
    email: email,
    password: password,
  })
    .then(checkStatus)
    .then(loginSucceed)
    .catch(loginFailed)
    .finally((_) => {
      submitBtn.classList.remove("is-loading");
    });
};

onDOMContentLoaded = (function () {
  err_msg_email = document.getElementById("err_msg_email");
  err_msg_password = document.getElementById("err_msg_password");
  email_field = document.getElementsByName("email")[0];
  password_field = document.getElementsByName("password")[0];
  submitBtn = document.querySelector("#submit_button");
  submitBtn.addEventListener("click", login);
})();
