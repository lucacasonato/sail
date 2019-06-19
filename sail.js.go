package main

//go:generate go run sail.js_gen.go
const sailJS = "(function() {\n    let oldonkeydown\n    function startReloadUI() {\n        const div = document.createElement(\"div\")\n        div.className = \"msgbox-overlay\"\n        div.style.opacity = 1\n        div.style.textAlign = \"center\"\n        div.innerHTML = `<div class=\"msgbox\">\n    <div class=\"msg\">Rebuilding container</div>\n    </div>`\n        // Prevent keypresses.\n        oldonkeydown = document.body.onkeydown\n        document.body.onkeydown = ev => {\n            ev.stopPropagation()\n        }\n        document.getElementsByClassName(\"monaco-workbench\")[0].appendChild(div)\n    }\n\n    function removeElementsByClass(className) {\n        let elements = document.getElementsByClassName(className);\n        for (let e of elements) {\n            e.parentNode.removeChild(e)\n        }\n    }\n\n    function stopReloadUI() {\n        document.body.onkeydown = oldonkeydown\n        removeElementsByClass(\"msgbox-overlay\")\n    }\n\n    let tty\n    let rebuilding\n    function rebuild() {\n        if (rebuilding) {\n            return\n        }\n        rebuilding = true\n\n        const tsrv = window.ide.workbench.terminalService\n\n        if (tty == null) {\n            tty = tsrv.createTerminal({\n                name: \"sail\",\n                isRendererOnly: true,\n            }, false)\n        } else {\n            tty.clear()\n        }\n        let oldTTY = tsrv.getActiveInstance()\n        tsrv.setActiveInstance(tty)\n        tsrv.showPanel(true)\n\n        startReloadUI()\n\n        const ws = new WebSocket(\"ws://\" + location.host + \"/sail/api/v1/reload\")\n        ws.onmessage = (ev) => {\n            const msg = JSON.parse(ev.data)\n            const out = atob(msg.v).replace(/\\n/g, \"\\n\\r\")\n            tty.write(out)\n        }\n        ws.onclose = (ev) => {\n            if (ev.code === 1000) {\n                tsrv.setActiveInstance(oldTTY)\n            } else {\n                alert(\"reload failed; please see logs in sail terminal\")\n            }\n            stopReloadUI()\n            rebuilding = false\n        }\n    }\n\n    window.addEventListener(\"ide-ready\", () => {\n        class rebuildAction extends window.ide.workbench.action {\n            run() {\n                rebuild()\n            }\n        }\n\n        window.ide.workbench.actionsRegistry.registerWorkbenchAction(new window.ide.workbench.syncActionDescriptor(rebuildAction, \"sail.rebuild\", \"Rebuild container\", {\n            primary: ((1 << 11) >>> 0) | 48 // That's cmd + R. See vscode source for the magic numbers.\n        }), \"sail: Rebuild container\", \"sail\");\n\n        const statusBarService = window.ide.workbench.statusbarService\n        statusBarService.addEntry({\n            text: \"rebuild\",\n            tooltip: \"Rebuild sail container\",\n            command: \"sail.rebuild\"\n        }, 0)\n    })\n}())\n"
