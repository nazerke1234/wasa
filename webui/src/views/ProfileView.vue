<template>
  <div class="profile-wrapper">
    <div class="toolbar">
      <h2>{{ displayName }}, welcome back</h2>
      <div class="toolbar-actions">
        <button @click="reloadProfile">Refresh</button>
        <button @click="signOut" class="logout-btn">Log Out</button>
        <button @click="createGroup" class="create-btn">New Group</button>
      </div>
    </div>

    <div class="profile-content">
      <div class="photo-section">
        <img v-if="photoSrc" :src="photoSrc" alt="Avatar" class="avatar" />
        <div v-else class="no-photo">No photo available</div>
      </div>
      <div class="info-section">
        <h3>{{ displayName }}</h3>
        <div class="field-row">
          <input v-model="newName" placeholder="Change username" />
          <button @click="changeName" :disabled="!newName || newName === displayName">Update</button>
        </div>
        <div class="field-row">
          <input type="file" accept="image/*" @change="onSelectPhoto" />
          <button @click="changePhoto" :disabled="!selectedPhoto">Upload</button>
        </div>
      </div>
    </div>

    <ErrorMsg v-if="error" :msg="error" />
  </div>
</template>

<script>
import axios from '../services/axios';
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
  name: 'UserProfileView',
  components: { ErrorMsg },
  data() {
    return {
      displayName: '',
      photoSrc: null,
      newName: '',
      selectedPhoto: null,
      error: null
    };
  },
  mounted() {
    this.loadProfile();
  },
  methods: {
    async loadProfile() {
      try {
        const token = localStorage.getItem('token');
        if (!token) return this.$router.push('/');

        const { data } = await axios.get('/users/photo', {
          headers: { Authorization: `Bearer ${token}` }
        });

        this.displayName = localStorage.getItem('name');
        this.photoSrc = data.photo ? `data:image/jpeg;base64,${data.photo}` : null;
      } catch (e) {
        console.error(e);
        this.error = 'Could not load profile.';
      }
    },
    onSelectPhoto(e) {
      const file = e.target.files[0];
      if (file) this.selectedPhoto = file;
    },
    async changePhoto() {
      try {
        const token = localStorage.getItem('token');
        const reader = new FileReader();
        reader.onload = async () => {
          const base64 = reader.result.split(',')[1];
          
          await axios.put('/users/photo', base64, {
            headers: { 
              Authorization: `Bearer ${token}`,
              'Content-Type': 'image/jpeg' // или 'image/png' в зависимости от файла
            }
          });
          
          alert('Photo updated!');
          this.loadProfile();
        };
        reader.readAsDataURL(this.selectedPhoto);
      } catch (e) {
        console.error(e);
        this.error = 'Could not update photo.';
      }
    }
    async changeName() {
      try {
        const token = localStorage.getItem('token');
        
        // Проверка на существующее имя
        if (this.newName === this.displayName) {
          this.error = "New username cannot be the same as current";
          return;
        }
    
        const response = await axios.put('/users/name', 
          { name: this.newName }, 
          {
            headers: { 
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json'
            }
          }
        );
    
        if (response.status === 200) {
          localStorage.setItem('name', response.data.name);
          this.displayName = response.data.name;
          this.newName = '';
          this.error = null;
          alert('Username updated successfully!');
        }
        
      } catch (error) {
        console.error('Update error:', error);
        
        if (error.response?.status === 400 || error.response?.status === 409) {
          this.error = "Username already exists or is invalid";
        } else {
          this.error = "Could not update username. Please try again.";
        }
      }
    }
    reloadProfile() {
      this.loadProfile();
    },
    signOut() {
      localStorage.clear();
      this.$router.push('/');
    },
    createGroup() {
      this.$router.push('/new-group');
    }
  }
};
</script>

<style scoped>
.profile-wrapper {
  max-width: 900px;
  margin: auto;
  padding: 2rem;
  background-color: #f9fafb;
  border-radius: 12px;
  box-shadow: 0 0 10px rgba(0,0,0,0.05);
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}
.toolbar-actions {
  display: flex;
  gap: 0.5rem;
}
button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  background-color: #007bff;
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
.create-btn {
  background-color: #28a745;
}
.create-btn:hover {
  background-color: #218838;
}
.profile-content {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
}
.photo-section {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.no-photo {
  font-size: 14px;
  color: #999;
}
.info-section {
  flex: 1;
}
.field-row {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}
input[type="text"], input[type="file"] {
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 14px;
  flex: 1;
  max-width: 300px;
}

</style>
