<template>
  <div class="group-info-wrapper">
    <div class="group-header">
      <div class="photo-box">
        <img v-if="groupPhoto" :src="groupPhoto" alt="Group Photo" class="group-img" />
      </div>
      <div class="info-box">
        <h1 class="group-title">{{ groupName }}</h1>
        <div class="edit-name">
          <input
            v-model="newGroupName"
            placeholder="New group name"
            maxlength="16"
            minlength="3"
          />
          <button @click="updateGroupName" :disabled="!newGroupName || newGroupName === groupName">
            Rename
          </button>
        </div>
        <div class="edit-photo">
          <input type="file" @change="handleGroupPhotoUpload" accept="image/*" />
          <button @click="updateGroupPhoto" :disabled="!newGroupPhoto">Change Photo</button>
        </div>

        <div class="member-section">
          <h3>Add Members</h3>
          <form @submit.prevent="searchUsers" class="search-form">
            <input v-model="query" placeholder="Search users..." />
            <button type="submit">Search</button>
          </form>

          <div v-if="showResults">
            <div v-if="users.length === 0">No results for "{{ lastQuery }}"</div>
            <div v-else class="user-cards">
              <div v-for="user in users" :key="user.id" class="user-card">
                <span>@{{ user.name }}</span>
                <button
                  @click="handleAddToGroup(user.id)"
                  :disabled="isMember(user.id)"
                >
                  {{ isMember(user.id) ? 'Member' : 'Add' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="leave-box">
          <button class="leave-btn" @click="leaveGroup">Leave Group</button>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<script>
import axios from '../services/axios';
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
  name: 'GroupEdit',
  components: { ErrorMsg },
  data() {
    return {
      token: localStorage.getItem('token'),
      groupId: this.$route.params.uuid,
      groupName: localStorage.getItem('groupName'),
      groupPhoto: null,
      newGroupName: '',
      newGroupPhoto: null,
      errormsg: null,
      query: '',
      lastQuery: '',
      users: [],
      loading: false,
      showResults: false,
      members: [],
    };
  },
  methods: {
    async fetchGroupProfile() {
      try {
        if (!this.token) return this.$router.push('/');
        const { data } = await axios.get(`/groups/${this.groupId}`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        this.groupPhoto = data.groupPhoto ? `data:image/*;base64,${data.groupPhoto}` : null;
        this.members = data.members;
      } catch (err) {
        console.error(err);
        this.errormsg = 'Could not load group data.';
      }
    },
    handleGroupPhotoUpload(e) {
      const file = e.target.files[0];
      if (file) this.newGroupPhoto = file;
    },
    async updateGroupPhoto() {
      if (!this.newGroupPhoto) return;
      try {
        const formData = new FormData();
        formData.append('photo', this.newGroupPhoto);
        await axios.put(`/groups/${this.groupId}/photo`, formData, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        alert('Group photo updated!');
        this.newGroupPhoto = null;
        this.fetchGroupProfile();
      } catch (err) {
        console.error(err);
        this.errormsg = 'Photo update failed.';
      }
    },
    async updateGroupName() {
      if (!this.newGroupName || this.newGroupName === this.groupName) return;
      try {
        await axios.put(
          `/groups/${this.groupId}/name`,
          { groupName: this.newGroupName },
          { headers: { Authorization: `Bearer ${this.token}` } }
        );
        alert('Group name updated!');
        localStorage.setItem('groupName', this.newGroupName);
        this.groupName = this.newGroupName;
        this.newGroupName = '';
      } catch (err) {
        console.error(err);
        this.errormsg = 'Name update failed.';
      }
    },
    async leaveGroup() {
      if (!confirm('Are you sure you want to leave?')) return;
      try {
        await axios.delete(`/groups/${this.groupId}`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        this.$router.push('/groups');
      } catch (err) {
        console.error(err);
        this.errormsg = 'Leave group failed.';
      }
    },
    async searchUsers() {
      if (!this.query.trim()) {
        this.errormsg = 'Enter a name.';
        this.showResults = false;
        return;
      }
      this.loading = true;
      this.errormsg = null;
      try {
        const { data } = await axios.get('/search', {
          params: { username: this.query },
        });
        this.users = data;
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (err) {
        console.error(err);
        this.errormsg = 'Search failed.';
      } finally {
        this.loading = false;
      }
    },
    isMember(id) {
      return this.members.includes(id);
    },
    async handleAddToGroup(id) {
      if (this.isMember(id)) return;
      try {
        await axios.post(
          `/groups/${this.groupId}`,
          { userId: id },
          { headers: { Authorization: `Bearer ${this.token}` } }
        );
        this.members.push(id);
      } catch (err) {
        console.error(err);
        this.errormsg = 'Add failed.';
      }
    },
  },
  mounted() {
    this.fetchGroupProfile();
  },
};
</script>

<style scoped>
.group-info-wrapper {
  padding: 2rem;
  max-width: 900px;
  margin: auto;
  background-color: #f9fafb;
  border-radius: 12px;
  box-shadow: 0 0 10px rgba(0,0,0,0.05);
}
.group-header {
  display: flex;
  flex-direction: row;
  gap: 1.5rem;
  align-items: flex-start;
}
.photo-box {
  flex: 0 0 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  background: #e0e0e0;
  display: flex;
  justify-content: center;
  align-items: center;
}
.group-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.info-box {
  flex: 1;
}
.group-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 1rem;
}
.edit-name,
.edit-photo,
.search-form,
.leave-box {
  margin-bottom: 1rem;
  display: flex;
  gap: 0.5rem;
}
input[type="text"],
input[type="file"],
button {
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 14px;
}
button {
  background: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}
button:hover:not(:disabled) {
  background: #0056b3;
}
button:disabled {
  background: #cccccc;
  cursor: not-allowed;
}
.member-section {
  margin-top: 2rem;
}
.user-cards {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 1rem;
}
.user-card {
  display: flex;
  justify-content: space-between;
  background: white;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.leave-btn {
  background-color: #dc3545;
  color: white;
}
.leave-btn:hover {
  background-color: #c82333;
}
@media (max-width: 768px) {
  .group-header {
    flex-direction: column;
    align-items: center;
  }
  .photo-box {
    width: 100px;
    height: 100px;
  }
  .group-title {
    text-align: center;
  }
}
</style>