package controller

type HandlerInterface interface {
	Add(obj interface{})
	Update(oldObj, newObj interface{})
	Delete(obj interface{})
}
