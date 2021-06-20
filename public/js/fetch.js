const fetchData = async (endpoint, { body, ...customConfig } = {}) => {
  const headers = { 'Content-Type': 'application/json' };
  const config = {
    method: customConfig.method || body ? 'POST' : 'GET',
    ...customConfig,
    headers: {
      ...headers,
      ...customConfig.headers,
    },
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  return fetch(endpoint, config);
};

/*
 * const myHeaders = new Headers({
 *   "Content-Type": "text/plain",
 *   "Content-Length": content.length.toString(),
 *   "X-Custom-Header": "ProcessThisImmediately",
 * });
 *
 * return fetch(endpoint, {
 *   method: method, // *GET, POST, PUT, DELETE, etc.
 *   mode: mode, // no-cors, *cors, same-origin
 *   cache: cache, // *default, no-cache, reload, force-cache, only-if-cached
 *   credentials: credentials, // include, *same-origin, omit
 *   headers: new Headers(header),
 *   redirect: redirect, // manual, *follow, error
 *   referrerPolicy: referrer, // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-endpoint
 *   body: body,
 * });
 */
