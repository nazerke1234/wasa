<template>
  <div class="groups-wrapper">
    <header class="toolbar">
      <h2>{{ username }}, your groups</h2>
      <div class="actions">
        <button @click="refreshGroups">Refresh</button>
        <button @click="newGroup">New Group</button>
        <button @click="logOut" class="logout-btn">Log Out</button>
      </div>
    </header>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <div v-if="loading" class="loading">Loading...</div>

    <div v-else-if="groups.length === 0" class="no-groups">
      No groups found.
    </div>

    <div v-else class="group-list">
      <div
        v-for="group in groups"
        :key="group.id"
        class="group-card"
        @click="goToGroup(group.id, group.name)"
      >
        <img
          v-if="group.conversationPhoto?.String"
          :src="'data:image/png;base64,' + group.conversationPhoto.String"
          alt="Group"
          class="group-image"
        />
        <div class="group-name">{{ group.name }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import ErrorMsg from '../components/ErrorMsg.vue';
import axios from '../services/axios';

export default {
  name: 'GroupsView',
  components: {
    ErrorMsg,
  },
  data() {
    return {
      username: localStorage.getItem('name') || 'Guest',
      groups: [],
      loading: false,
      errormsg: null,
    };
  },
  methods: {
    async refreshGroups() {
      this.loading = true;
      this.errormsg = null;

      try {
        const token = localStorage.getItem('token');
        if (!token) return this.$router.push('/');

        const { data } = await axios.get('/groups', {
          headers: { Authorization: `Bearer ${token}` },
        });

        this.groups = data || [];
      } catch (err) {
        console.error('Error loading groups:', err);
        this.errormsg = 'Failed to load groups. Please try again.';
      } finally {
        this.loading = false;
      }
    },
    goToGroup(id, name) {
      localStorage.setItem('groupName', name);
      this.$router.push(`/groups/${id}`);
    },
    newGroup() {
      this.$router.push('/new-group');
    },
    logOut() {
      this.$router.push('/');
    },
  },
  mounted() {
    this.refreshGroups();
  },
};
</script>

<style scoped>
.groups-wrapper {
  max-width: 800px;
  margin: auto;
  padding: 2rem;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 2rem;
}

.toolbar h2 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 1rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 14px;
  border: none;
  background: #007bff;
  color: white;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.logout-btn {
  background-color: #6c757d;
}

.logout-btn:hover {
  background-color: #5a6268;
}

.loading,
.no-groups {
  text-align: center;
  color: #666;
  margin-top: 2rem;
}

.group-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.group-card {
  background: #fff;
  padding: 1rem;
  border-radius: 10px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: transform 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.group-card:hover {
  transform: translateY(-3px);
}

.group-image {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 0.5rem;
}

.group-name {
  font-weight: 600;
  font-size: 16px;
}
</style>