const fetchData = async (url = "", method = "GET", mode = "cors", cache = "no-cache", credentials = "same-origin", data = {}, header = {}, redirect = "follow", referrer = "no-referrer", body = "", isJSON = false) => {
  // const myHeaders = new Headers({
  //   "Content-Type": "text/plain",
  //   "Content-Length": content.length.toString(),
  //   "X-Custom-Header": "ProcessThisImmediately",
  // });

  if (method === "GET" || method === "HEAD") {
    body = null;
  } else if (isJSON) {
    body = JSON.stringify(body);
  }
  return fetch(url, {
    method: method, // *GET, POST, PUT, DELETE, etc.
    mode: mode, // no-cors, *cors, same-origin
    cache: cache, // *default, no-cache, reload, force-cache, only-if-cached
    credentials: credentials, // include, *same-origin, omit
    headers: new Headers(header),
    redirect: redirect, // manual, *follow, error
    referrerPolicy: referrer, // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: body,
  });
};
