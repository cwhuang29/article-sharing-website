const sendResetPasswordEndpoint = '/password/reset';
const errMsg = {
  empty: "This field can't be empty.",
  short: 'Password must be at least 8 characters long.',
  notMatch: "The password confirm and password didn't match.",
};

const getInputValue = () => {
  return {
    password: inputPassword.value.trim(),
    passwordConfirm: inputPasswordConfirm.value.trim(),
  };
};

const validateInput = ({ password, passwordConfirm } = values) => {
  let canSubmit = true;

  if (password.length == 0) {
    canSubmit = false;
    err_msg_password.innerText = errMsg.empty;
  } else if (password.length < 8) {
    canSubmit = false;
    err_msg_password.innerText = errMsg.short;
  } else {
    err_msg_password.innerText = '';
  }

  if (passwordConfirm.length == 0) {
    canSubmit = false;
    err_msg_passwordConfirm.innerText = errMsg.empty;
  } else if (passwordConfirm.length < 8) {
    canSubmit = false;
    err_msg_passwordConfirm.innerText = errMsg.short;
  } else if (password != passwordConfirm) {
    canSubmit = false;
    err_msg_passwordConfirm.innerText = errMsg.notMatch;
  } else {
    err_msg_passwordConfirm.innerText = '';
  }

  return canSubmit;
};

const clearInputBox = () => {
  inputPassword.value = '';
  inputPasswordConfirm.value = '';
};

const sendResetPasswordEmailSucceed = async (resp) => {
  resp.json().then(function (data) {
    showNoticeMsg(data.msgHead, data.msgBody);
  });
  clearInputBox();
  setTimeout(() => (window.location.href = resp.headers.get('Location')), 2000);
  return Promise.resolve();
};

const sendResetPasswordEmailFailed = async (resp) => {
  resp.json().then(function (data) {
    showErrMsg(data.errHead, data.errBody);
  });
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  } else {
    return Promise.resolve(resp);
  }
};

const sendResetPasswordForm = () => {
  submitBtn.classList.add('is-loading');

  const values = getInputValue();
  const ok = validateInput(values);

  if (!ok) {
    submitBtn.classList.remove('is-loading');
    return;
  }

  /*
   * In some frameworks' design, they extract these two values from the HTML from
   * The email handling is the same as I did, and the token is stored in a hidden input <input type="hidden" name="token" value="{{ $token }}">
   * Since I send requests by fetching instead of submitting form, I'll just take their values from the URL
   */
  const token = window.location.href.split('/').pop().split('?')[0]; // The structure of URL is "/reset/password/<token>?email=<email>"
  const email = new URLSearchParams(window.location.search).get('email');
  const allValues = {
    email: email,
    password: values.password,
    token: token,
  };

  const csrfToken = document.querySelector('meta[name="csrf-token"]').content;
  const headers = { 'X-CSRF-TOKEN': csrfToken };

  fetchData(sendResetPasswordEndpoint, { body: allValues, method: 'PUT', headers: headers })
    .then(checkStatus)
    .then(sendResetPasswordEmailSucceed)
    .catch(sendResetPasswordEmailFailed)
    .finally((_) => {
      submitBtn.classList.remove('is-loading');
    });
};

const passwordResetFormHandler = () => {
  err_msg_password = document.getElementById('err_msg_password');
  err_msg_passwordConfirm = document.getElementById('err_msg_password_confirm');
  inputPassword = document.getElementsByName('password')[0];
  inputPasswordConfirm = document.getElementsByName('password_confirm')[0];
  submitBtn = document.querySelector('#submitButton');

  submitBtn.addEventListener('click', sendResetPasswordForm);
};

onDOMContentLoaded = (function () {
  passwordResetFormHandler();
})();
