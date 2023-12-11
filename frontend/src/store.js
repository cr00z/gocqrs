import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueNativeSock from 'vue-native-websocket';

const BACKEND_URL = 'https://localhost:8080';
const PUSHER_URL = 'ws://localhost:8080/pusher';

const SET_MESSAGES = 'SET_MESSAGES';
const CREATE_MESSAGE = 'CREATE_MESSAGE';
const SEARCH_SUCCESS = 'SEARCH_SUCCESS';
const SEARCH_ERROR = 'SEARCH_ERROR';

const MESSAGE_CREATED = 1;

Vue.use(Vuex);

const store = new Vuex.Store({
    state: {
        messages: [],
        searchResults: [],
    },
    mutations: {
        SOCKET_ONOPEN(state, event) {},
        SOCKET_ONCLOSE(state, event) {},
        SOCKET_ONERROR(state, event) {
            console.error(event);
        },
        SOCKET_ONMESSAGE(state, message) {
            switch (message.kind) {
                case MESSAGE_CREATED:
                    this.commit(CREATE_MESSAGE, { id: message.id, body: message.body });
            }
        },
        [SET_MESSAGES](state, messages) {
            state.messages = messages;
        },
        [CREATE_MESSAGE](state, message) {
            state.messages = [message, ...state.messages];
        },
        [SEARCH_SUCCESS](state, messages) {
            state.searchResults = messages;
        },
        [SEARCH_ERROR](state) {
            state.searchResults = [];
        },
    },
    actions: {
        getMessages({ commit }) {
            axios
                .get(`${BACKEND_URL}/messages`)
                .then(({ data }) => {
                    commit(SET_MESSAGES, data);
                })
                .catch((err) => console.error(err));
        },
        async createMessage({ commit }, message) {
            const { data } = await axios.post(`${BACKEND_URL}/messages`, null, {
                params: {
                    body: message.body,
                },
            });
        },
        async searchMessages({ commit }, query) {
            if (query.length == 0) {
                commit(SEARCH_SUCCESS, []);
                return;
            }
            axios
                .get(`${BACKEND_URL}/search`, {
                    params: { query },
                })
                .then(({ data }) => commit(SEARCH_SUCCESS, data))
                .catch((err) => {
                    console.error(err);
                    commit(SEARCH_ERROR);
                });
        },
    },
});

Vue.use(VueNativeSock, PUSHER_URL, { store, format: 'json' });

store.dispatch('getMessages');

export default store;