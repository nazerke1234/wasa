<template>
  <div class="message-interface">
    <div class="message-header">
      <div class="header-image" v-if="chatImage">
        <img :src="'data:image/jpeg;base64,' + chatImage" alt="Chat Image" />
      </div>
      <h3>{{ chatTitle }}</h3>
    </div>
    <div class="message-list" ref="messageList">
      <p v-if="messageList.length === 0">No messages available...</p>
      <div
        v-for="msg in messageList"
        :key="msg.id"
        class="message-item"
        :class="msg.senderId === currentUserToken ? 'own' : 'others'"
        :style="msg.senderId !== currentUserToken && chatCategory === 'group' ? { paddingLeft: '45px' } : {}"
      >
        <div v-if="chatCategory === 'group' && msg.senderId !== currentUserToken" class="user-avatar">
          <img :src="'data:image/jpeg;base64,' + msg.senderPhoto" alt="User Avatar" />
        </div>
        <div class="message-body">
          <div v-if="msg.replyTo" class="response-preview">
            <small>Response to {{ msg.replySenderName || 'Unknown' }}: {{ msg.replyContent }}</small>
            <img
              v-if="msg.replyAttachment"
              :src="'data:image/jpeg;base64,' + msg.replyAttachment"
              alt="Response Attachment"
              class="response-image"
            />
          </div>
          <p v-if="msg.content.startsWith('<strong>Forwarded')" v-html="msg.content"></p>
          <p v-else>
            <strong>
              {{ msg.senderId === currentUserToken ? 'You' : (msg.senderName || 'Unknown User') }}:
            </strong>
            {{ msg.content }}
          </p>
          <div v-if="msg.attachment" class="media-container">
            <img :src="'data:image/jpeg;base64,' + msg.attachment" alt="Media" class="media-display" />
          </div>
          <small>{{ formatTime(msg.timestamp) }}</small>
          <div v-if="msg.reactionCount > 0" class="emoji-count">
            ❤️ × {{ msg.reactionCount }}
            <div class="emoji-users">
              <ul>
                <li v-for="(user, index) in msg.reactingUserNames" :key="user">
                  {{ index + 1 }}. {{ user }}
                </li>
              </ul>
            </div>
          </div>
          <div class="interaction-buttons">
            <button v-if="msg.senderId !== currentUserToken" class="interaction-btn response-btn" @click.stop="setupResponse(msg)">
              ↩
            </button>
            <button
              v-if="msg.senderId !== currentUserToken"
              class="interaction-btn emoji-btn"
              :class="{ 'active-emoji': (msg.reactingUserNames || []).includes(currentUserName) }"
              :disabled="msg.reactionLoading"
              @click.stop="handleEmoji(msg)"
            >
              ❤️
            </button>
            <button class="interaction-btn share-btn" @click.stop="displayShareMenu(msg.id)">
              →
            </button>
            <button v-if="msg.senderId === currentUserToken" class="interaction-btn remove-btn" @click.stop="removeMessage(msg)">
              ✖
            </button>
          </div>
          <div v-if="messageSettings[msg.id]?.displayShareMenu" class="share-menu" @click.stop>
            <select id="share-select" class="share-select" v-model="messageSettings[msg.id].chosenChatId">
              <option value="" disabled>Choose chat</option>
              <option v-for="chat in messageSettings[msg.id].availableChats" :key="chat.id" :value="chat.id">
                {{ chat.name }}
              </option>
              <option value="new">New contact</option>
            </select>
            <div v-if="messageSettings[msg.id].chosenChatId === 'new'" class="find-contact">
              <input type="text" v-model="messageSettings[msg.id].contactSearch" placeholder="Enter contact name" @input="findContact(msg.id)" />
              <ul v-if="messageSettings[msg.id].contactMatches.length > 0" class="contact-matches">
                <li v-for="contact in messageSettings[msg.id].contactMatches" :key="contact.id" @click="pickContact(contact, msg.id)" class="contact-match">
                  {{ contact.name }}
                </li>
              </ul>
            </div>
            <div class="share-actions">
              <button class="action-btn" v-if="messageSettings[msg.id].chosenChatId !== 'new'" :disabled="!messageSettings[msg.id].chosenChatId" @click.stop="shareMessage(messageSettings[msg.id].chosenChatId, msg.id)">
                Share
              </button>
              <button class="action-btn" v-if="messageSettings[msg.id].chosenChatId === 'new'" :disabled="!messageSettings[msg.id].pickedContactId" @click.stop="shareToContact(messageSettings[msg.id].pickedContactId, msg.id)">
                Share
              </button>
              <button class="action-btn" @click.stop="hideShareMenu(msg.id)">Close</button>
            </div>
            <div v-if="messageSettings[msg.id].availableChats.length === 0">
              No chats available.
            </div>
          </div>
        </div>
        <div class="message-state" v-if="msg.status && msg.senderId !== currentUserToken">
          {{ msg.status }}
        </div>
      </div>
    </div>
    <div v-if="responseMessage" class="response-preview-box">
      <div class="response-info">
        <strong>Responding to {{ responseMessage.senderName || 'Unknown' }}:</strong>
        <span class="response-text">{{ responseMessage.content }}</span>
        <img v-if="responseMessage.attachment" :src="'data:image/jpeg;base64,' + responseMessage.attachment" alt="Response Media" class="response-media-preview" />
      </div>
      <button class="cancel-response-btn" @click="clearResponse">✖</button>
    </div>
    <div class="message-composer">
      <input type="file" ref="mediaInput" style="display: none" accept="image/*, .gif" @change="processMediaSelection" />
      <button class="media-attach-btn" @click="activateMediaInput">
        Add Media or GIF
        <span v-if="chosenMedia" class="media-indicator">🖼️</span>
      </button>
      <input v-model="newMessage" class="composer-input" type="text" placeholder="Write your message..." @input="updateSendButton" />
      <button v-if="newMessage.trim() || chosenMedia" class="submit-btn" @click="submitMessage">
        Submit
      </button>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
export default {
  name: "MessageInterface",
  data() {
    return {
      newMessage: "",
      messageList: [],
      currentUserToken: localStorage.getItem("token"),
      chatTitle: localStorage.getItem("conversationName") || "Unknown Chat",
      chatImage: null,
      chatCategory: null,
      chatId: this.$route.params.uuid,
      messageSettings: {},
      chosenMedia: null,
      pollingTimer: null,
      initialLoad: true,
      responseMessage: null
    };
  },
  computed: {
    currentUserName() {
      return localStorage.getItem("name");
    }
  },
  methods: {
    activateMediaInput() {
      this.$refs.mediaInput.click();
    },
    processMediaSelection(event) {
      this.chosenMedia = event.target.files[0];
    },
    async submitMessage() {
      const authToken = localStorage.getItem("token");
      if (!authToken) {
        this.$router.push({ path: "/" });
        return;
      }
      const dataForm = new FormData();
      dataForm.append("content", this.newMessage);
      if (this.responseMessage) {
        dataForm.append("replyTo", this.responseMessage.id);
      }
      if (this.chosenMedia) {
        dataForm.append("attachment", this.chosenMedia);
      }
      await axios.post(`/conversations/${this.chatId}/message`, dataForm, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      this.newMessage = "";
      this.chosenMedia = null;
      this.$refs.mediaInput.value = "";
      this.responseMessage = null;
      await this.loadMessages();
      this.$nextTick(() => {
        this.scrollToEnd();
      });
    },
    async loadMessages() {
      const authToken = localStorage.getItem("token");
      if (!authToken) {
        this.$router.push({ path: "/" });
        return;
      }
      const result = await axios.get(`/conversations/${this.chatId}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      this.messageList = (result.data.messages || []).map(message => ({
        ...message,
        reactingUserNames: message.reactingUserNames || [],
        showReactionUsers: false
      }));
      if (result.data.name) {
        this.chatTitle = result.data.name;
      }
      if (result.data.conversationPhoto && result.data.conversationPhoto.String) {
        this.chatImage = result.data.conversationPhoto.String;
      } else {
        this.chatImage = null;
      }
      this.chatCategory = result.data.type || "direct";
      this.$nextTick(() => {
        if (this.initialLoad) {
          this.scrollToEnd();
          this.initialLoad = false;
        }
      });
    },
    scrollToEnd() {
      const messageArea = this.$refs.messageList;
      if (messageArea) {
        messageArea.scrollTop = messageArea.scrollHeight;
      }
    },
    async handleEmoji(message) {
      const authToken = localStorage.getItem("token");
      if (!authToken || message.senderId === this.currentUserToken) return;
      const userReacted = (message.reactingUserNames || []).includes(this.currentUserName);
      try {
        if (userReacted) {
          await axios.delete(`/conversations/${this.chatId}/message/${message.id}/comment`, {
            headers: { Authorization: `Bearer ${authToken}` }
          });
        } else {
          await axios.post(`/conversations/${this.chatId}/message/${message.id}/comment`, {},
            { headers: { Authorization: `Bearer ${authToken}` } }
          );
        }
      } catch (error) {
        console.error("Error handling emoji", error);
      } finally {
        await this.loadMessages();
      }
    },
    async removeMessage(message) {
      const authToken = localStorage.getItem("token");
      if (!authToken) {
        this.$router.push({ path: "/" });
        return;
      }
      await axios.delete(`/conversations/${this.chatId}/message/${message.id}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      this.messageList = this.messageList.filter(m => m.id !== message.id);
    },
    formatTime(timestamp) {
      const dateObj = new Date(timestamp);
      return dateObj.toLocaleString();
    },
    displayShareMenu(messageId) {
      this.hideAllMenus();
      if (!this.messageSettings[messageId]) {
        this.messageSettings[messageId] = {
          displayShareMenu: true,
          availableChats: [],
          chosenChatId: "",
          contactSearch: "",
          contactMatches: [],
          pickedContactId: ""
        };
        this.loadAvailableChats(messageId);
      } else {
        this.messageSettings[messageId].displayShareMenu = !this.messageSettings[messageId].displayShareMenu;
      }
    },
    hideShareMenu(messageId) {
      if (this.messageSettings[messageId]) {
        this.messageSettings[messageId].displayShareMenu = false;
      }
    },
    hideAllMenus() {
      for (const id in this.messageSettings) {
        this.messageSettings[id].displayShareMenu = false;
      }
    },
    handleExternalClick(event) {
      const messageBody = this.$el.querySelector('.message-body');
      if (messageBody && !messageBody.contains(event.target)) {
        this.hideAllMenus();
      }
    },
    async loadAvailableChats(messageId) {
      const authToken = localStorage.getItem("token");
      const result = await axios.get('/conversations', {
        headers: { Authorization: `Bearer ${authToken}` }
      });
      const chats = result.data.filter(chat => chat.id !== this.chatId);
      this.messageSettings[messageId].availableChats = chats;
    },
    async findContact(messageId) {
      const searchTerm = this.messageSettings[messageId].contactSearch;
      if (!searchTerm.trim()) {
        this.messageSettings[messageId].contactMatches = [];
        return;
      }
      const authToken = localStorage.getItem("token");
      const result = await axios.get('/search', {
        params: { username: searchTerm },
        headers: { Authorization: `Bearer ${authToken}` }
      });
      this.messageSettings[messageId].contactMatches = result.data;
    },
    pickContact(contact, messageId) {
      this.messageSettings[messageId].pickedContactId = contact.id;
      this.messageSettings[messageId].contactSearch = contact.name;
      this.messageSettings[messageId].contactMatches = [];
    },
    async shareToContact(selectedContactId, messageId) {
      const authToken = localStorage.getItem("token");
      if (!authToken) {
        this.$router.push({ path: "/" });
        return;
      }
      const chatResult = await axios.post(
        `/conversations`,
        { senderId: authToken, recipientId: selectedContactId },
        { headers: { Authorization: `Bearer ${authToken}` } }
      );
      const destinationChatId = chatResult.data.conversationId;
      const sharerName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.chatId}/message/${messageId}/forward`,
        { sourceMessageId: messageId, targetConversationId: destinationChatId, forwarderName: sharerName },
        { headers: { Authorization: `Bearer ${authToken}` } }
      );
      alert("Message shared successfully!");
      this.hideShareMenu(messageId);
    },
    async shareMessage(targetChatId, messageId) {
      const message = this.messageList.find(m => m.id === messageId);
      if (!message) return;
      const authToken = localStorage.getItem("token");
      const sharerName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.chatId}/message/${messageId}/forward`,
        { targetConversationId: targetChatId, forwarderName: sharerName },
        { headers: { Authorization: `Bearer ${authToken}` } }
      );
      alert("Message shared successfully!");
      this.hideShareMenu(messageId);
    },
    setupResponse(message) {
      this.responseMessage = message;
    },
    clearResponse() {
      this.responseMessage = null;
    },
    updateSendButton() {
      // Method for future send button state updates
    }
  },
  mounted() {
    this.loadMessages();
    this.pollingTimer = setInterval(() => {
      this.loadMessages();
    }, 5000);
    document.addEventListener("click", this.handleExternalClick);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.handleExternalClick);
    clearInterval(this.pollingTimer);
  }
};
</script>

<style scoped>
.message-interface {
  display: flex;
  flex-direction: column;
  height: 92vh;
  overflow: hidden;
}
.message-header {
  display: flex;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}
.header-image {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}
.header-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
.message-list {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  border-top: 1px solid #ccc;
  border-bottom: 1px solid #ccc;
}
.message-item {
  position: relative;
  max-width: 70%;
  margin-bottom: 10px;
  border-radius: 10px;
  padding: 10px;
  background-color: #e0f2f1;
}
.message-item.own {
  margin-left: auto;
  background-color: #d1e7dd;
}
.user-avatar {
  position: absolute;
  left: 10px;
  top: 10px;
  width: 30px;
  height: 30px;
}
.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
.message-body {
  position: relative;
  min-height: 40px;
}
.message-item p {
  margin: 0;
  color: #333;
  word-wrap: break-word;
  white-space: pre-wrap;
}
.message-item small {
  margin-top: 5px;
  color: #666;
  font-size: 0.8em;
}
.media-container {
  margin-top: 8px;
  width: 300px;
  height: 300px;
  overflow: hidden;
  border: 1px solid #ddd;
  border-radius: 8px;
}
.media-display {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.interaction-buttons {
  position: absolute;
  top: 0;
  right: -50px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.message-item.own .interaction-buttons {
  left: -50px;
  right: auto;
}
.interaction-btn {
  position: static;
  width: 17px;
  height: 17px;
  border: 1px solid #aaa;
  border-radius: 50%;
  background-color: white;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
  opacity: 0.9;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  padding: 0;
}
.interaction-btn:hover {
  opacity: 1;
}
.response-btn {
  font-size: 10px;
  margin-right: 5px;
}
.share-menu {
  position: absolute;
  top: 30px;
  right: 0;
  background-color: #ffffff;
  border-radius: 5px;
  padding: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  z-index: 100;
  width: 250px;
}
.share-select {
  width: 100%;
  padding: 8px;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  margin-bottom: 8px;
  font-size: 14px;
}
.share-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 10px;
}
.action-btn {
  background-color: #128c7e;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}
.action-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
.find-contact input {
  width: 100%;
  padding: 6px;
  margin-bottom: 4px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.contact-matches {
  list-style: none;
  padding: 0;
  margin: 0 0 6px 0;
  max-height: 100px;
  overflow-y: auto;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.contact-match {
  padding: 4px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}
.contact-match:hover {
  background-color: #f0f0f0;
}
.media-indicator {
  font-size: 18px;
  margin-left: 5px;
}
.message-state {
  position: absolute;
  bottom: 5px;
  right: 10px;
  font-size: 12px;
  color: #555;
}
.emoji-users ul {
  margin: 0;
  padding: 0;
  list-style: none;
  font-size: 0.8em;
  color: #444;
}
.emoji-users li {
  margin: 2px 0;
}
.response-preview-box {
  background-color: #f0f0f0;
  border-left: 4px solid #128c7e;
  padding: 8px;
  margin: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.response-info {
  font-size: 0.9em;
  color: #444;
}
.response-image,
.response-media-preview {
  width: 21px;
  height: 21px;
  object-fit: cover;
  margin-left: 10px;
  border-radius: 4px;
}
.cancel-response-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #888;
}
.message-composer {
  display: flex;
  align-items: flex-start;
  flex-wrap: wrap;
  padding: 10px;
  gap: 8px;
}
.media-attach-btn {
  background-color: #25d366;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  margin-right: 10px;
  font-size: 14px;
  padding: 10px 15px;
}
.media-attach-btn:hover {
  background-color: #20b358;
}
.composer-input {
  flex: 1;
  min-width: 200px;
  padding: 12px;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  font-size: 14px;
  outline: none;
}
.submit-btn {
  background-color: #128c7e;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 20px;
  margin-left: 10px;
  cursor: pointer;
  font-size: 14px;
}
.submit-btn:hover {
  background-color: #0f7c6a;
}
@media (max-width: 600px) {
  .conversation-block p {
    -webkit-line-clamp: 3;
    line-clamp: 3;
  }
}
</style>
