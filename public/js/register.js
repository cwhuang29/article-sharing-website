const errInputMsg = {
  empty: "This field can't be empty.",
  passwordInconsistency: 'Password and confirm password does not match.',
  passwordTooShort: 'Passwords must be at least 8 characters long.',
  emailFormatInvalid: 'The email format is not correct.',
};
const registerEndpoint = '/register';

const setErrMsgBanner = (ele, msg) => {
  ele.innerText = msg;
};

const isValidEmail = (e) => {
  const filter = /^\s*[\w\-\+_]+(\.[\w\-\+_]+)*\@[\w\-\+_]+\.[\w\-\+_]+(\.[\w\-\+_]+)*\s*$/;
  return String(e).search(filter) != -1;
};

const getInputValue = () => {
  let g = [...gender].filter((g) => g.checked).map((g) => g.value)[0] || '';
  return {
    firstName: firstName.value.trim(),
    lastName: lastName.value.trim(),
    password: password.value.trim(),
    passwordConfirm: passwordConfirm.value.trim(),
    email: email.value.trim(),
    gender: g,
    major: major.value,
  };
};

const validateInput = (values) => {
  let isValid = true;

  for (const [_, val] of Object.entries(values)) {
    if (val.length == 0) {
      isValid = false;
    }
  }

  const { firstName, lastName, password, passwordConfirm, email, gender, major } = values;

  if (firstName.length == 0) {
    setErrMsgBanner(firstNameErr, errInputMsg.empty);
  } else {
    setErrMsgBanner(firstNameErr, '');
  }

  if (lastName.length == 0) {
    setErrMsgBanner(lastNameErr, errInputMsg.empty);
  } else {
    setErrMsgBanner(lastNameErr, '');
  }

  if (password.length == 0) {
    setErrMsgBanner(passwordErr, errInputMsg.empty);
  } else if (password.length < 8) {
    isValid = false;
    setErrMsgBanner(passwordErr, errInputMsg.passwordTooShort);
  } else {
    setErrMsgBanner(passwordErr, '');
  }

  if (passwordConfirm.length == 0) {
    setErrMsgBanner(passwordConfirm, errInputMsg.empty);
  } else if (passwordConfirm.length < 8) {
    isValid = false;
    setErrMsgBanner(passwordConfirm, errInputMsg.passwordTooShort);
  } else {
    setErrMsgBanner(passwordConfirm, '');
  }

  if (password.length > 7 && passwordConfirm.length > 7) {
    if (password != passwordConfirm) {
      isValid = false;
      setErrMsgBanner(passwordConfirmErr, errInputMsg.passwordInconsistency);
    } else {
      setErrMsgBanner(passwordErr, '');
      setErrMsgBanner(passwordConfirmErr, '');
    }
  }

  if (email.length == 0) {
    setErrMsgBanner(emailErr, errInputMsg.empty);
  } else if (!isValidEmail(email)) {
    isValid = false;
    setErrMsgBanner(emailErr, errInputMsg.emailFormatInvalid);
  } else {
    setErrMsgBanner(emailErr, '');
  }

  if (gender.length == 0) {
    setErrMsgBanner(genderErr, errInputMsg.empty);
  } else {
    setErrMsgBanner(genderErr, '');
  }

  if (major.length == 0) {
    setErrMsgBanner(majorErr, errInputMsg.empty);
  } else {
    setErrMsgBanner(majorErr, '');
  }

  return isValid;
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  } else {
    return Promise.resolve(resp);
  }
};

const creationSucceed = async (resp) => {
  window.location.href = resp.headers.get('Location');
  return Promise.resolve();
};

const creationFailed = async (resp) => {
  if (resp.status == 500) {
    resp.json().then(function (data) {
      showErrMsg(data.errHead, data.errBody);
    });
  } else if (resp.status == 400) {
    resp.json().then(function (data) {
      if (data.bindingError) {
        c(data.errHead);
        showErrMsg('An Error Occurred !', 'Please reload the page and try again.');
      } else {
        for (var key in data.errTags) {
          document.getElementById(`err_msg_${key}`).innerText = data.errTags[key];
        }
      }
    });
  } else if (resp.status == 409) {
    resp.json().then((data) => showErrMsg(data.errHead, data.errBody));
  } else {
    showErrMsg('An Error Occurred !', 'Please reload the page and try again.');
  }
};

const submitNewPost = () => {
  submitBtn.classList.add('is-loading');

  let values = getInputValue();
  let res = validateInput(values);

  if (!res) {
    submitBtn.classList.remove('is-loading');
    return;
  }

  const transformedValues = {
    first_name: values.firstName,
    last_name: values.lastName,
    password: values.password,
    email: values.email,
    gender: values.gender,
    major: values.major,
  };

  fetchData(registerEndpoint, {
    body: transformedValues,
    redirect: 'follow',
  })
    .then(checkStatus)
    .then(creationSucceed)
    .catch(creationFailed)
    .finally((_) => {
      submitBtn.classList.remove('is-loading');
    });
};

const registerHandler = () => {
  firstName = document.getElementsByName('first_name')[0];
  lastName = document.getElementsByName('last_name')[0];
  password = document.getElementsByName('password')[0];
  passwordConfirm = document.getElementsByName('password_confirm')[0];
  email = document.getElementsByName('email')[0];
  gender = document.getElementsByName('gender'); // Array
  major = document.getElementsByName('major')[0];

  firstNameErr = document.getElementById('err_msg_first_name');
  lastNameErr = document.getElementById('err_msg_last_name');
  passwordErr = document.getElementById('err_msg_password');
  passwordConfirmErr = document.getElementById('err_msg_password_confirm');
  emailErr = document.getElementById('err_msg_email');
  genderErr = document.getElementById('err_msg_gender');
  majorErr = document.getElementById('err_msg_major');

  submitBtn = document.querySelector('#submit_button');
  submitBtn.addEventListener('click', submitNewPost);
};

onDOMContentLoaded = (function () {
  registerHandler();
})();
