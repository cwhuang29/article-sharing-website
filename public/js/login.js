const loginEndpoint = '/login';

const validateEmail = (email) => {
  const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
};

const getInputValue = () => {
  return {
    email: email_field.value.trim(),
    password: password_field.value.trim(),
  };
};

const validateInput = ({ email, password } = values) => {
  let canSubmit = true;

  if (email.length == 0) {
    canSubmit = false;
    err_msg_email.innerText = "The email can't be empty.";
  } else if (!validateEmail(email)) {
    canSubmit = false;
    err_msg_email.innerText = 'Please fill in the correct email.';
  } else {
    err_msg_email.innerText = '';
  }

  if (password.length == 0) {
    canSubmit = false;
    err_msg_password.innerText = "The password can't be empty.";
  } else {
    err_msg_password.innerText = '';
  }

  return canSubmit;
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  } else {
    return Promise.resolve(resp);
  }
};

const loginSucceed = async (resp) => {
  // Copy from overview.js (to make sure admin users can see admin only articles once login since data will not be fetched if session storage is not empty)
  window.sessionStorage.removeItem('offset');
  window.sessionStorage.removeItem('overviewContent');
  window.location.href = resp.headers.get('Location');
  return Promise.resolve();
};

const loginFailed = async (resp) => {
  if (resp.status == 500) {
    resp.json().then(function (data) {
      showErrMsg(data.errHead, data.errBody);
    });
  } else {
    resp.json().then(function (data) {
      showErrMsg(data.errHead, data.errBody);
      for (var key in data.errTags) {
        document.getElementById(`err_msg_${key}`).innerText = data.err[key];
      }
    });
  }
};

const login = async () => {
  submitBtn.classList.add('is-loading');

  const values = getInputValue();
  const ok = validateInput(values);

  if (!ok) {
    submitBtn.classList.remove('is-loading');
    return;
  }

  fetchData(loginEndpoint, { body: values })
    .then(checkStatus)
    .then(loginSucceed)
    .catch(loginFailed)
    .finally((_) => submitBtn.classList.remove('is-loading'));
};

const loginHandler = () => {
  err_msg_email = document.getElementById('err_msg_email');
  err_msg_password = document.getElementById('err_msg_password');
  email_field = document.getElementsByName('email')[0];
  password_field = document.getElementsByName('password')[0];
  submitBtn = document.querySelector('#submit_button');
  submitBtn.addEventListener('click', login);
};

onDOMContentLoaded = (function () {
  loginHandler();
})();
