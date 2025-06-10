ls ./backend/go/*go  | xargs -I {} cat {} > backend_full.txt
ls ./frontend/HTML/*.html | xargs -I {} cat {} > frontend_html_full.txt
ls ./frontend/CSS/*.css | xargs -I {} cat {} > frontend_css_full.txt
ls ./frontend/JS/*.js | xargs -I {} cat {} > frontend_js_full.txt