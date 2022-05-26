import http from 'k6/http';
import exec from 'k6/execution';
import { Httpx } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { describe } from 'https://jslib.k6.io/expect/0.0.5/index.js';

export let options = {
	vus: 100,
	duration: '10s',
};

let session = new Httpx({
	baseURL: 'http://app:8080',
	headers: {
		'Authorization': 'whatever'
	}
})

export default function () {
	let key = exec.vu.idInTest;
	let value = 'content';
	describe('set-get-del', (t) => {
		let res = session.get(`/v1/${key}`);
		t.expect(res.status).as(`miss`).toEqual(404);

		res = session.post(`/v1/${key}/${value}`);
		t.expect(res.status).as(`set`).toEqual(201);

		res = session.get(`/v1/${key}`);
		t.expect(res.status).as(`hit`).toEqual(200);

		res = session.delete(`/v1/${key}`);
		t.expect(res.status).as(`del`).toEqual(200);
	});
}
