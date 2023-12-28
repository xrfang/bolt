import { defineConfig } from 'vite';
import { createHtmlPlugin }  from 'vite-plugin-html';
import fs from 'fs';
import { resolve } from 'path';

function readHtmlFiles(dir) {
    const files = fs.readdirSync(dir);
    const htmlContents = {};
    files.forEach(file => {
        if (file.endsWith('.html')) {
            const fullPath = resolve(dir, file);
            const content = fs.readFileSync(fullPath, 'utf-8');
            const key = file.replace('.html', '');
            htmlContents[key] = content;
        }
    });

    return htmlContents;
}

const htmlParts = readHtmlFiles('shared/');

export default defineConfig({
  plugins: [
    createHtmlPlugin({
      inject: {
        data: htmlParts
      },
      minify: true,
    }),
  ],
  // ... 其他 Vite 配置
});