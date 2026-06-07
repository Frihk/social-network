const fs = require('fs');
const path = require('path');

// Simple validation script to check frontend structure
const checks = [
  { name: 'Auth context', path: 'src/context/AuthContext.jsx' },
  { name: 'Auth library', path: 'src/lib/auth.js' },
  { name: 'Middleware', path: 'src/middleware.js' },
  { name: 'Root layout', path: 'src/app/layout.jsx' },
  { name: 'Home page', path: 'src/app/page.jsx' },
  { name: 'Login page', path: 'src/app/login/page.jsx' },
  { name: 'Register page', path: 'src/app/register/page.jsx' },
  { name: 'Profile page', path: 'src/app/profile/[id]/page.jsx' },
  { name: 'Global styles', path: 'src/styles/globals.css' },
  { name: 'Package config', path: 'package.json' },
  { name: 'Next config', path: 'next.config.js' },
  { name: 'Tailwind config', path: 'tailwind.config.js' },
];

console.log('🔍 Validating Frontend Structure...\n');

let allValid = true;
checks.forEach(check => {
  const fullPath = path.join(__dirname, check.path);
  if (fs.existsSync(fullPath)) {
    console.log(`✅ ${check.name}: ${check.path}`);
  } else {
    console.log(`❌ ${check.name}: ${check.path} (NOT FOUND)`);
    allValid = false;
  }
});

console.log();
if (allValid) {
  console.log('✅ All files present!\n');
  console.log('📝 Next Steps:');
  console.log('1. Run: npm install');
  console.log('2. Run: npm run dev');
  console.log('3. Open: http://localhost:3000');
  console.log('4. Test registration and login flows');
} else {
  console.log('❌ Some files are missing!');
  process.exit(1);
}
