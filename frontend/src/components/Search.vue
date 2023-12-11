<template>
  <div>
    <input @keyup="searchMessages" v-model.trim="query" type="text" class="form-control" placeholder="Search...">
    <div class="mt-4">
      <Message v-for="message in messages" :key="message.id" :message="message" />
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import Message from '@/components/Message';

export default {
  data() {
    return {
      query: '',
    };
  },
  computed: mapState({
    messages: (state) => state.searchResults,
  }),
  methods: {
    searchMessages() {
      if (this.query != this.lastQuery) {
        this.$store.dispatch('searchMessages', this.query);
        this.lastQuery = this.query;
      }
    },
  },
  components: {
    Message,
  },
};
</script>