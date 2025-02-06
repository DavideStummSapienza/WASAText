import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

// Interceptor, to append automaticaly the `identifier` to all requests
instance.interceptors.request.use((config) => {
	const identifier = sessionStorage.getItem("identifier");
	if (identifier) {
	  config.headers.Authorization = `Bearer ${identifier}`;
	}
	return config;
  });

export default instance;
