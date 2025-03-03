export default class Helper {
    static getFile(key, item) {
        return "/api/files/" + item.collectionId + "/" + item.id + "/" + item[key];
    }
}