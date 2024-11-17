// Supports ES6
// import { create, Whatsapp } from '@wppconnect-team/wppconnect';
const wppconnect = require('@wppconnect-team/wppconnect');

wppconnect
  .create({
    phoneNumber: '6287834100533',
    catchLinkCode: (str) => console.log('Code: ' + str),
  })
  .then((client) => start(client))
  .catch((error) => console.log(error));
