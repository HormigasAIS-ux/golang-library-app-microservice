const { execSync } = require('child_process');

function getCommitAuthor() {
  return execSync('git log -1 --pretty=format:"%an"').toString().trim();
}

const commitAuthor = getCommitAuthor();

if (commitAuthor === "Zakky" || commitAuthor === "zakku116") {
  console.log("Error: Zakky no tiene permisos para hacer commit.");
  process.exit(1);  // Esto previene el commit
} else {
  process.exit(0);  // Permite el commit
}