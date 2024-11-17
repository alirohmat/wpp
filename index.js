const wppconnect = require('@wppconnect-team/wppconnect');
const { exec } = require('child_process');

// Daftar ID WhatsApp yang diizinkan (owner)
const owners = ['6282132789470@s.whatsapp.net']; // Gantilah dengan ID WhatsApp Anda

function start(client) {
  client.onMessage(async (message) => {
    // Cek apakah pesan diawali dengan $ dan pengirim adalah owner
    if (message.body.startsWith('$') && owners.includes(message.from)) {
      try {
        await client.sendText(message.from, 'Processing...');

        // Mengambil perintah shell dari pesan setelah simbol "$"
        const command = message.body.slice(1).trim(); // Hapus "$" dan spasi tambahan

        // Mengeksekusi perintah shell
        exec(command, (err, stdout, stderr) => {
          if (err) {
            // Mengirimkan pesan error jika ada masalah dalam eksekusi
            return client.sendText(message.from, `Error: ${err.message}\nExit Code: ${err.code}\nSignal: ${err.signal}\n${stderr}`);
          }
          if (stdout) {
            // Mengirimkan hasil output standar jika ada
            return client.sendText(message.from, `Output:\n${stdout}`);
          }
          if (stderr) {
            // Mengirimkan pesan jika ada output error
            return client.sendText(message.from, `stderr: ${stderr}`);
          }
          // Jika tidak ada output, kirim pesan konfirmasi
          client.sendText(message.from, 'Command executed successfully with no output.');
        });
      } catch (e) {
        // Mengirimkan pesan error jika terjadi kesalahan saat mengeksekusi perintah
        client.sendText(message.from, `An error occurred: ${e.toString()}`);
      }
    } else {
      // Jika bukan owner yang mengirim pesan
      client.sendText(message.from, 'You are not authorized to use this command.');
    }
  });
}

wppconnect
  .create({
    phoneNumber: '6287834100533', // Nomor WhatsApp bot
    catchLinkCode: (str) => console.log('Code: ' + str),
  })
  .then((client) => start(client))
  .catch((error) => console.log(error));
