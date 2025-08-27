<template>
  <div class="chat-box">
    <header class="chat-header">
      <img v-if="chatPhoto" :src="'data:image/jpeg;base64,' + chatPhoto" class="chat-avatar" />
      <h2>{{ chatTitle }}</h2>
    </header>

    <section class="chat-content" ref="messageArea">
      <p v-if="chatMessages.length === 0">No messages available.</p>
      <div
        v-for="msg in chatMessages"
        :key="msg.id"
        :class="['chat-bubble', msg.senderId === myToken ? 'me' : 'they']"
      >
        <div v-if="isGroup && msg.senderId !== myToken" class="user-photo">
          <img :src="'data:image/jpeg;base64,' + msg.senderPhoto" alt="user" />
        </div>

        <div class="bubble-body">
          <div v-if="msg.replyTo" class="reply-wrapper">
            <span class="reply-author">Reply to {{ msg.replySenderName || 'Unknown' }}</span>
            <span class="reply-text">: {{ msg.replyContent }}</span>
            <img
              v-if="msg.replyAttachment"
              :src="'data:image/jpeg;base64,' + msg.replyAttachment"
              class="reply-img"
            />
          </div>

          <p v-if="msg.content.startsWith('<strong>Forwarded')" v-html="msg.content"></p>
          <p v-else><strong>{{ msg.senderId === myToken ? 'You' : msg.senderName }}</strong>: {{ msg.content }}</p>

          <img
            v-if="msg.attachment"
            :src="'data:image/jpeg;base64,' + msg.attachment"
            class="attached-img"
          />

          <small class="timestamp">{{ formatDate(msg.timestamp) }}</small>

          <div class="reactions" v-if="msg.reactionCount > 0">
            ❤️ x {{ msg.reactionCount }}
            <ul class="react-list">
              <li v-for="(name, i) in msg.reactingUserNames" :key="name">{{ i + 1 }}. {{ name }}</li>
            </ul>
          </div>

          <div class="bubble-actions">
            <button v-if="msg.senderId !== myToken" @click.stop="startReply(msg)">↩</button>
            <button
              v-if="msg.senderId !== myToken"
              :class="{ reacted: msg.reactingUserNames.includes(myName) }"
              @click.stop="toggleHeart(msg)"
            >❤️</button>
            <button @click.stop="forwardUI(msg.id)">→</button>
            <button v-if="msg.senderId === myToken" @click.stop="removeMessage(msg)">✖</button>
          </div>

          <div
            v-if="forwardMenus[msg.id]?.visible"
            class="forward-ui"
            @click.stop
          >
            <select v-model="forwardMenus[msg.id].targetChatId">
              <option disabled value="">Select chat</option>
              <option v-for="c in forwardMenus[msg.id].available" :key="c.id" :value="c.id">{{ c.name }}</option>
              <option value="new">New contact</option>
            </select>

            <div v-if="forwardMenus[msg.id].targetChatId === 'new'">
              <input v-model="forwardMenus[msg.id].query" placeholder="Username" @input="findUser(msg.id)" />
              <ul v-if="forwardMenus[msg.id].results.length">
                <li
                  v-for="u in forwardMenus[msg.id].results"
                  :key="u.id"
                  @click="chooseUser(msg.id, u)"
                >{{ u.name }}</li>
              </ul>
            </div>

            <div class="forward-buttons">
              <button
                v-if="forwardMenus[msg.id].targetChatId !== 'new'"
                @click.stop="forwardToChat(forwardMenus[msg.id].targetChatId, msg.id)"
              >Send</button>
              <button
                v-else
                :disabled="!forwardMenus[msg.id].chosenUserId"
                @click.stop="forwardToUser(forwardMenus[msg.id].chosenUserId, msg.id)"
              >Send</button>
              <button @click.stop="closeForward(msg.id)">Cancel</button>
            </div>

            <div v-if="!forwardMenus[msg.id].available.length">No chats.</div>
          </div>
        </div>
      </div>
    </section>

    <div v-if="pendingReply" class="reply-bar">
      <span>Replying to {{ pendingReply.senderName || 'Unknown' }}: {{ pendingReply.content }}</span>
      <button @click="cancelReply">✖</button>
    </div>

    <footer class="chat-input-area">
      <input ref="fileInput" type="file" hidden accept="image/*,.gif" @change="fileChosen" />
      <button class="file-btn" @click="triggerFile">📎</button>
      <input
        v-model="messageText"
        class="msg-field"
        placeholder="Type message..."
        @input="handleTyping"
      />
      <button class="send-btn" @click="submitMessage" :disabled="!messageText.trim() && !selectedFile">Send</button>
    </footer>
  </div>
</template>

<script>
import axios from "../services/axios";

export default {
  name: "ChatRoom",
  data() {
    return {
      messageText: "",
      chatMessages: [],
      myToken: localStorage.getItem("token"),
      chatTitle: localStorage.getItem("conversationName") || "Chat",
      chatPhoto: null,
      isGroup: false,
      convoId: this.$route.params.uuid,
      forwardMenus: {},
      selectedFile: null,
      pollingId: null,
      firstInit: true,
      pendingReply: null,
    };
  },
  computed: {
    myName() {
      return localStorage.getItem("name") || "Me";
    },
  },
  methods: {
    triggerFile() {
      this.$refs.fileInput.click();
    },
    fileChosen(e) {
      this.selectedFile = e.target.files[0];
    },
    async submitMessage() {
      const token = this.myToken;
      if (!token) return this.$router.push("/");

      const form = new FormData();
      form.append("content", this.messageText);
      if (this.pendingReply) form.append("replyTo", this.pendingReply.id);
      if (this.selectedFile) form.append("attachment", this.selectedFile);

      await axios.post(`/chats/${this.convoId}/message`, form, {
        headers: { Authorization: `Bearer ${token}` },
      });

      this.messageText = "";
      this.selectedFile = null;
      this.pendingReply = null;
      this.$refs.fileInput.value = "";
      await this.fetchAll();
      this.$nextTick(() => this.scrollToEnd());
    },
    async fetchAll() {
      const token = this.myToken;
      if (!token) return;

      const res = await axios.get(`/chats/${this.convoId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data = res.data;
      this.chatMessages = (data.messages || []).map(m => ({
        ...m,
        reactingUserNames: m.reactingUserNames || [],
      }));

      this.chatTitle = data.name || this.chatTitle;
      this.chatPhoto = data.conversationPhoto?.String || null;
      this.isGroup = data.type === "group";

      if (this.firstInit) {
        this.$nextTick(() => this.scrollToEnd());
        this.firstInit = false;
      }
    },
    scrollToEnd() {
      const box = this.$refs.messageArea;
      if (box) box.scrollTop = box.scrollHeight;
    },
    formatDate(ts) {
      return new Date(ts).toLocaleString();
    },
    async toggleHeart(msg) {
      const token = this.myToken;
      if (!token || msg.senderId === token) return;

      const reacted = (msg.reactingUserNames || []).includes(this.myName);
      const endpoint = `/chats/${this.convoId}/message/${msg.id}/comment`;

      try {
        if (reacted) {
          await axios.delete(endpoint, { headers: { Authorization: `Bearer ${token}` } });
        } else {
          await axios.post(endpoint, {}, { headers: { Authorization: `Bearer ${token}` } });
        }
      } catch (e) {
        console.warn("Reaction toggle failed", e);
      } finally {
        await this.fetchAll();
      }
    },
    async removeMessage(msg) {
      const token = this.myToken;
      if (!token) return;

      await axios.delete(`/chats/${this.convoId}/message/${msg.id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      this.chatMessages = this.chatMessages.filter(m => m.id !== msg.id);
    },
    forwardUI(id) {
      this.hideAllMenus();

      if (!this.forwardMenus[id]) {
        this.forwardMenus[id] = {
          visible: true,
          available: [],
          targetChatId: "",
          query: "",
          results: [],
          chosenUserId: "",
        };
        this.loadChats(id);
      } else {
        this.forwardMenus[id].visible = !this.forwardMenus[id].visible;
      }
    },
    hideAllMenus() {
      Object.values(this.forwardMenus).forEach(m => (m.visible = false));
    },
    closeForward(id) {
      if (this.forwardMenus[id]) this.forwardMenus[id].visible = false;
    },
    async loadChats(id) {
      const token = this.myToken;
      const res = await axios.get("/chats", {
        headers: { Authorization: `Bearer ${token}` },
      });
      this.forwardMenus[id].available = res.data.filter(c => c.id !== this.convoId);
    },
    async findUser(id) {
      const q = this.forwardMenus[id].query;
      if (!q.trim()) return (this.forwardMenus[id].results = []);

      const token = this.myToken;
      const res = await axios.get("/search", {
        headers: { Authorization: `Bearer ${token}` },
        params: { username: q },
      });
      this.forwardMenus[id].results = res.data;
    },
    chooseUser(id, u) {
      const menu = this.forwardMenus[id];
      menu.chosenUserId = u.id;
      menu.query = u.name;
      menu.results = [];
    },
    async forwardToUser(uid, mid) {
      const token = this.myToken;
      if (!token) return;

      const convRes = await axios.post(
        "/chats",
        { senderId: token, recipientId: uid },
        { headers: { Authorization: `Bearer ${token}` } }
      );

      const newConvId = convRes.data.conversationId;
      await this.forwardMessage(mid, newConvId);
    },
    async forwardToChat(cid, mid) {
      await this.forwardMessage(mid, cid);
    },
    async forwardMessage(mid, targetId) {
      const token = this.myToken;
      const name = this.myName;

      await axios.post(
        `/chats/${this.convoId}/message/${mid}/forward`,
        { targetChatId: targetId, forwarderName: name },
        { headers: { Authorization: `Bearer ${token}` } }
      );

      alert("Message forwarded!");
      this.closeForward(mid);
    },
    startReply(msg) {
      this.pendingReply = msg;
    },
    cancelReply() {
      this.pendingReply = null;
    },
    handleTyping() {},
  },
  mounted() {
    this.fetchAll();
    this.pollingId = setInterval(this.fetchAll, 5000);
    document.addEventListener("click", this.hideAllMenus);
  },
  beforeUnmount() {
    clearInterval(this.pollingId);
    document.removeEventListener("click", this.hideAllMenus);
  },
};
</script>

<style scoped>
.chat-box {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #eef2f5;
  border-bottom: 1px solid #ccc;
}
.chat-avatar {
  width: 40px;
  height: 40px;
  margin-right: 10px;
  border-radius: 50%;
  object-fit: cover;
}
.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  background: #f9f9f9;
}
.chat-bubble {
  max-width: 70%;
  margin: 8px 0;
  padding: 10px;
  border-radius: 10px;
  background-color: #dceefb;
  position: relative;
}
.chat-bubble.me {
  align-self: flex-end;
  background-color: #c8e6c9;
}
.user-photo img {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 8px;
}
.bubble-body {
  display: flex;
  flex-direction: column;
}
.attached-img {
  width: 200px;
  height: auto;
  margin-top: 5px;
  border-radius: 5px;
  object-fit: cover;
}
.timestamp {
  font-size: 0.75rem;
  color: #777;
  margin-top: 4px;
}
.bubble-actions button {
  background: none;
  border: none;
  font-size: 12px;
  cursor: pointer;
  margin-right: 4px;
}
.forward-ui {
  background: white;
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 6px;
  margin-top: 6px;
}
.forward-buttons {
  display: flex;
  gap: 8px;
  margin-top: 6px;
}
.reply-bar {
  background: #f1f1f1;
  padding: 8px;
  border-left: 3px solid #007bff;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reply-bar button {
  border: none;
  background: transparent;
  font-size: 18px;
  cursor: pointer;
}
.chat-input-area {
  display: flex;
  align-items: center;
  padding: 10px;
  border-top: 1px solid #ccc;
}
.file-btn {
  font-size: 18px;
  border: none;
  background: none;
  cursor: pointer;
}
.msg-field {
  flex: 1;
  margin: 0 10px;
  padding: 10px;
  border-radius: 20px;
  border: 1px solid #ddd;
}
.send-btn {
  background: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
}
.send-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
.react-list {
  font-size: 0.75rem;
  color: #444;
  margin-top: 4px;
  list-style: none;
  padding-left: 0;
}
</style>


