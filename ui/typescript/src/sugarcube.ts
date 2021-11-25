//-----------------------------------------------------------------------------
// URL utilities
//-----------------------------------------------------------------------------

export function getPath() {
    return window.location.pathname;
}

export function parseURLPattern(pattern) {
    // remove trailing slash in pattern
    if (pattern[pattern.length - 1] == '/') {
        pattern = pattern.substr(0, pattern.length - 1);
    }

    // explode and scan for placeholders.
    let placeholders = {};
    let tokens = pattern.split('/');
    for (let i = 0; i < tokens.length; i++) {
        let token = tokens[i];
        if (token[0] != '<' || token[token.length - 1] != '>') {
            continue;
        }

        token = token.substring(1, token.length - 1);
        let parts = token.split(':');

        if (parts.length != 2) {
            // invalid pattern format.
            return null;
        } else if (parts[1].length <= 0) {
            // placeholder name is empty.
            return null;
        } else if (placeholders.hasOwnProperty(parts[1]) == true) {
            // placeholder name already exists.
            return null;
        } else if (parts[0] != 'str' && parts[0] != 'int') {
            // invalid type.
            return null;
        }

        placeholders[parts[1]] = ''
        tokens[i] = { 'type' : parts[0], 'name' : parts[1] };
    }

    return tokens;
}

export function renderPath(pattern, params) {
    let parsedPtn = parseURLPattern(pattern);
    if (parsedPtn == null) {
        return null;
    }

    let ret = '';
    for (let i = 0; i < parsedPtn.length; i++) {
        switch (typeof(parsedPtn[i])) {
            case 'string':
                ret = ret + parsedPtn[i];
                break;
            case 'object':
                let name = parsedPtn[i]['name'];
                let type = parsedPtn[i]['type'];

                if (params.hasOwnProperty(name) == false) {
                    return null;
                }

                let param = params[name];
                if (type == 'str') {
                    if (typeof(param) != 'string') {
                        return null;
                    }
                    ret = ret + param;
                } else if (type == 'int') {
                    if (typeof(param) != 'number') {
                        return null;
                    }
                    ret = ret + param.toString();
                } else {
                    return null;
                }
                break;
            default:
                return null;
        }

        if (i != parsedPtn.length - 1) {
            ret = ret + '/';
        }
    }

    return ret;
}

export function parsePath(pattern, path) {
    let parsedPtn = parseURLPattern(pattern);
    if (parsedPtn == null) {
        return null;
    }

    // remove trailing slash in url.
    if (pattern[pattern.length - 1] == '/') {
        path = path.substr(0, path.length - 1);
    }

    let pathTokens = path.split('/');
    if (pathTokens.length != parsedPtn.length) {
        return null;
    }

    // process pattern.
    let parsedValues = {};
    for (let i = 0; i < parsedPtn.length; i++) {
        switch (typeof(parsedPtn[i])) {
            case 'string':
                if (parsedPtn[i] != pathTokens[i]) {
                    return null;
                }
                break;
            case 'object':
                let name = parsedPtn[i]['name'];
                let type = parsedPtn[i]['type'];

                if (type == 'str') {
                    parsedValues[name] = pathTokens[i];
                } else if (type == 'int') {
                    let value = parseInt(pathTokens[i]);
                    if (value == NaN) {
                        return null;
                    }
                    parsedValues[name] = value;
                } else {
                    return null;
                }
                break;
            default:
                return null;
        }
    }

    return parsedValues;
}

//-----------------------------------------------------------------------------
// AJAX calls
//-----------------------------------------------------------------------------

export function ajaxGet(URL, callbackFunc) {
    let request = new XMLHttpRequest();

    if (callbackFunc != null) {
        request.onload = function() {
            callbackFunc(request.status, request.response);
        }
    }

    request.open('GET', URL);
    request.send();
}

export function ajaxDelete(URL, callbackFunc) {
    let request = new XMLHttpRequest();

    if (callbackFunc != null) {
        request.onload = function() {
            callbackFunc(request.status, request.response);
        }
    }

    request.open('DELETE', URL);
    request.send();
}

export function ajaxPost(URL, payload, callbackFunc) {
    let request = new XMLHttpRequest();

    if (callbackFunc != null) {
        request.onload = function() {
            callbackFunc(request.status, request.response);
        }
    }

    let toSend = null;
    if (payload != null) {
        switch (typeof(payload)) {
            case 'string':
                toSend = payload;
                break;
            case 'object':
                toSend = new FormData();
                for (let key in payload) {
                    toSend.append(key, payload[key])
                }
                break;
            default:
                return;
        }
    }

    request.open('POST', URL);
    request.send(toSend);
}

//-----------------------------------------------------------------------------
// CSS/JavaScript insertion
//-----------------------------------------------------------------------------

export function loadScript(scriptURL, onLoadFunc) {
    let scriptElem = document.createElement("script");
    scriptElem.src = scriptURL;
    scriptElem.onload = onLoadFunc;

    document.body.appendChild(scriptElem);
}

export function loadCSS(cssURL) {
    let cssElem = document.createElement("link");
    cssElem.rel = "stylesheet";
    cssElem.type = "text/css";
    cssElem.href = cssURL;

    document.head.appendChild(cssElem);
}

//-----------------------------------------------------------------------------
// Modal Popup
//-----------------------------------------------------------------------------

export const MODAL_BACKGROUND_ID = "sugarCubeModalBackground";
export const MODAL_BACKGROUND_STYLE = "position:fixed; padding:0; margin:0; top:0; left:0; width:100%; height:100%; background:rgba(180,180,180,0.5); backdrop-filter: blur(3px); z-index:65535;";
export const MODAL_DIALOG_STYLE = "position:fixed; top:50%; left:50%; background-color:#FFF; border-radius:10px; z-index:65536; padding: 40px; "
export const MODAL_CLOSE_STYLE = "font-size:24pt; position:absolute; top:0px; right:15px;"

export function modalPopup(width, height, contentHTML) {
    if (document.getElementById(MODAL_BACKGROUND_ID) != null) {
        return;
    }

    // background
    let backgroundElem = document.createElement("div");
    backgroundElem.id = MODAL_BACKGROUND_ID;
    backgroundElem.setAttribute('style', MODAL_BACKGROUND_STYLE);

    // close button of dialog
    let closeElem = document.createElement("a");
    closeElem.href = "#";
    closeElem.onclick = modalClose;
    closeElem.setAttribute('style', MODAL_CLOSE_STYLE);
    closeElem.innerText = "×";

    // dialog
    let dialogStyle = MODAL_DIALOG_STYLE.concat("width:", width.toString(), "px; height:", height.toString(), "px; margin-top:", (height / -2).toString(), "px; margin-left:", (width / -2).toString(), "px;");
    let dialogElem = document.createElement("div");
    dialogElem.setAttribute('style', dialogStyle)
    dialogElem.appendChild(closeElem);

    if (contentHTML != null) {
        let contentElem = document.createElement("div");
        contentElem.innerHTML = contentHTML;
        dialogElem.appendChild(contentElem)
    }

    backgroundElem.appendChild(dialogElem);
    document.body.appendChild(backgroundElem)
}

export function modalClose() {
    let backgroundElem = document.getElementById(MODAL_BACKGROUND_ID);
    if (backgroundElem != null) {
        document.body.removeChild(backgroundElem);
    }
}
