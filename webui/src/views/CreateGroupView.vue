<template>
  <div class="create-group-wrapper">
    <h1 class="title">New Group</h1>

    <div class="input-section">
      <label for="name">Group Name</label>
      <input
        id="name"
        v-model="groupName"
        class="text-input"
        type="text"
        placeholder="Type a group name"
      />
    </div>

    <form @submit.prevent="searchUsers" class="search-bar">
      <input
        v-model="query"
        class="search-input"
        type="text"
        placeholder="Search username..."
      />
      <button class="btn search-btn">Search</button>
    </form>

    <div v-if="error" class="alert">{{ error }}</div>
    <LoadingSpinner v-if="loading" />

    <div v-if="!loading && showResults" class="results">
      <h2 class="section-title">Results</h2>
      <div v-if="users.length" class="user-list">
        <div v-for="user in users" :key="user.id" class="user-item">
          <span>@{{ user.name }}</span>
          <button
            class="btn add-btn"
            @click="addUser(user)"
            :disabled="isSelected(user)"
          >
            {{ isSelected(user) ? 'Added' : 'Add' }}
          </button>
        </div>
      </div>
      <p v-else>No results for "{{ lastQuery }}"</p>
    </div>

    <div v-if="selectedUsers.length" class="selected-box">
      <h2 class="section-title">Members</h2>
      <div class="tags">
        <div v-for="user in selectedUsers" :key="user.id" class="tag">
          @{{ user.name }}
          <button class="remove-tag" @click="removeUser(user)">✕</button>
        </div>
      </div>
    </div>

    <div class="input-section">
      <label for="img">Group Image</label>
      <input id="img" ref="fileInput" type="file" accept="image/*" @change="handleImage" />
      <img v-if="preview" :src="preview" class="image-preview" />
    </div>

    <button class="btn create-btn" @click="createGroup" :disabled="!isValid">
      Create
    </button>
  </div>
</template>

<script>
import axios from '../services/axios';
import LoadingSpinner from '../components/LoadingSpinner.vue';

export default {
  name: 'GroupCreateView',
  components: { LoadingSpinner },
  data() {
    return {
      groupName: '',
      query: '',
      lastQuery: '',
      users: [],
      selectedUsers: [],
      file: null,
      preview: null,
      error: '',
      loading: false,
      showResults: false,
    };
  },
  computed: {
    isValid() {
      return this.groupName.trim() && this.selectedUsers.length && this.file;
    },
  },
  methods: {
    async searchUsers() {
      if (!this.query.trim()) {
        this.error = 'Enter a name to search.';
        this.showResults = false;
        return;
      }
      this.error = '';
      this.loading = true;
      this.users = [];
      try {
        const { data } = await axios.get('/search', {
          params: { username: this.query },
        });
        const me = localStorage.getItem('token');
        this.users = data.filter(u => String(u.id) !== me);
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (e) {
        this.error = 'Search failed. Try again.';
      } finally {
        this.loading = false;
      }
    },
    addUser(user) {
      if (!this.isSelected(user)) this.selectedUsers.push(user);
    },
    removeUser(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id);
    },
    isSelected(user) {
      return this.selectedUsers.some(u => u.id === user.id);
    },
    handleImage(e) {
      const file = e.target.files[0];
      if (file) {
        this.file = file;
        const reader = new FileReader();
        reader.onload = ev => (this.preview = ev.target.result);
        reader.readAsDataURL(file);
      }
    },
    async createGroup() {
      if (!this.isValid) return;
      const form = new FormData();
      form.append('name', this.groupName);
      form.append('image', this.file);
      const me = localStorage.getItem('token');
      const memberIds = [...this.selectedUsers.map(u => u.id), me];
      form.append('membersJson', JSON.stringify(memberIds));
      try {
        await axios.post('/groups', form, {
          headers: { 'Content-Type': 'multipart/form-data' },
        });
        alert('Group created');
        this.$router.push('/home');
      } catch (e) {
        this.error = 'Group creation failed.';
      }
    },
  },
};
</script>

<style scoped>
.create-group-wrapper {
  max-width: 640px;
  margin: 0 auto;
  padding: 20px;
  background: #fff;
}
.title {
  font-size: 26px;
  margin-bottom: 20px;
  text-align: center;
}
.input-section {
  margin-bottom: 20px;
  text-align: left;
}
.text-input, .search-input {
  width: 100%;
  padding: 10px;
  margin-top: 5px;
  border-radius: 6px;
  border: 1px solid #ccc;
}
.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}
.btn {
  padding: 8px 16px;
  border-radius: 5px;
  border: none;
  cursor: pointer;
}
.search-btn {
  background-color: #007bff;
  color: white;
}
.add-btn {
  background-color: #28a745;
  color: white;
}
.create-btn {
  background-color: #007bff;
  color: white;
  width: 100%;
  margin-top: 20px;
}
.alert {
  background: #f8d7da;
  color: #842029;
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 10px;
}
.results, .selected-box {
  margin-top: 20px;
}
.section-title {
  font-size: 18px;
  margin-bottom: 10px;
}
.user-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border: 1px solid #eee;
  border-radius: 5px;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.tag {
  background: #e2e3e5;
  padding: 5px 10px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  gap: 5px;
}
.remove-tag {
  border: none;
  background: transparent;
  cursor: pointer;
}
.image-preview {
  max-width: 100%;
  max-height: 200px;
  margin-top: 10px;
  border-radius: 8px;
}
</style>

