import { createI18n } from 'vue-i18n'

const messages = {
    en: {
        proxyTitle: 'Reverse Proxy',
        proxyDescription: 'List of proxy rules',
        firewallTitle: 'Firewall',
        firewallDescription: 'List of firewall rules',
        addmore: 'Add more',
        //////////////
        firewallBlockLog: 'Firewall block log',
        errorLog: 'Error log',
        forwardIp: 'forward-ip',
        originIp: 'origin-ip',
        error: 'error',
        time: 'time',
        host: 'host',
        //////////////
        packet: 'Packet',
        allow: 'Allow',
        block: 'Block',
        //////////////
        hostname: 'Hostname',
        firewall: 'Firewall',
        endpoint: 'Endpoint',
        actions: 'Actions',
        rulename: 'Rule name',
        //////////////
        issues: 'Issues',
    },
    kr: {
        proxyTitle: '리버스 프록시',
        proxyDescription: '프록시 설정 목록',
        firewallTitle: '방화벽',
        firewallDescription: '방화벽 규칙 목록',
        addmore: '규칙추가',
        //////////////
        firewallBlockLog: '방화벽 차단 로그',
        errorLog: '오류 로그',
        forwardIp: '전달IP',
        originIp: '원본IP',
        error: '오류',
        time: '시간',
        host: '호스트',
        //////////////
        packet: '패킷',
        allow: '허용',
        block: '차단',
        //////////////
        hostname: '호스트 이름',
        firewall: '방화벽 규칙',
        endpoint: '엔드포인트',
        actions: '작업',
        rulename: '규칙 이름',
        //////////////
        issues: '이슈',
    },
}
const i18n = createI18n({
    legacy: false,
    locale: 'en',
    fallbackLocale: 'en',
    messages,
})

export default i18n