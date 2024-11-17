const wppconnect = require('@wppconnect-team/wppconnect');
const { exec } = require('child_process');

// Daftar ID WhatsApp yang diizinkan (owner)
const owners = ['6282132789470@s.whatsapp.net']; // Gantilah dengan ID WhatsApp Anda

function start(client) {
  client.onMessage(async (message) => {
    // Cek apakah pesan diawali dengan $ dan pengirim adalah owner
    if (message.body.startsWith('$') && owners.includes(message.from)) {
      try {
        // Mengirimkan pesan "Processing..." untuk memberi tahu pengguna bahwa perintah sedang diproses
        await client.sendText(message.from, 'Processing...');

        // Mengambil perintah shell dari pesan setelah simbol "$"
        const command = message.body.slice(1).trim(); // Hapus "$" dan spasi tambahan

        // Jika perintah kosong, beri tahu pengguna
        if (!command) {
          return client.sendText(message.from, 'No command provided. Please provide a valid shell command.');
        }

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

  // Handle session status changes
  client.on('statusFind', (statusSession, session) => {
    console.log('Status Session: ', statusSession);
    console.log('Session name: ', session);
    if (statusSession === 'notLogged' || statusSession === 'qrReadFail' || statusSession === 'desconnectedMobile' || statusSession === 'qrReadError' || statusSession === 'autocloseCalled' || statusSession === 'browserClose') {
      console.log('Session is not logged in, QR code read failed, or disconnected from mobile. Re-creating session.');
      wppconnect.create({ phoneNumber: '6287834100533', catchLinkCode: (str) => console.log('Code: ' + str) })
        .then((newClient) => start(newClient))
        .catch((error) => console.log(error));
    }
  });

  // Handle errors
  client.on('error', (error) => {
    console.error('Error: ', error);
    if (error.message.includes('Execution context was destroyed')) {
      console.log('Execution context was destroyed. Re-creating session.');
      wppconnect.create({ phoneNumber: '6287834100533', catchLinkCode: (str) => console.log('Code: ' + str) })
        .then((newClient) => start(newClient))
        .catch((error) => console.log(error));
    }
  });
}

wppconnect
  .create({
    phoneNumber: '6287834100533', // Nomor WhatsApp bot
    catchLinkCode: (str) => console.log('Code: ' + str),
    statusFind: (statusSession, session) => {
      console.log('Status Session: ', statusSession);
      console.log('Session name: ', session);
      if (statusSession === 'notLogged' || statusSession === 'qrReadFail' || statusSession === 'desconnectedMobile' || statusSession === 'qrReadError' || statusSession === 'autocloseCalled' || statusSession === 'browserClose') {
        console.log('Session is not logged in, QR code read failed, or disconnected from mobile. Re-creating session.');
        wppconnect.create({ phoneNumber: '6287834100533', catchLinkCode: (str) => console.log('Code: ' + str) })
          .then((newClient) => start(newClient))
          .catch((error) => console.log(error));
      }
    },
    headless: true, // Headless chrome
    devtools: false, // Open devtools by default
    useChrome: true, // If false will use Chromium instance
    debug: false, // Opens a debug session
    logQR: true, // Logs QR automatically in terminal
    browserWS: '', // If you want to use browserWSEndpoint
    browserArgs: [''], // Parameters to be added into the chrome browser instance
    puppeteerOptions: {}, // Will be passed to puppeteer
  })
  .then((client) => start(client))
  .catch((error) => console.log(error));
