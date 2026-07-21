package service

import (
	"context"
	"fmt"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/integration/insales"
)

type Service struct {
	manager *core.AccountManager
}

func NewService(mgr *core.AccountManager) *Service {
	return &Service{manager: mgr}
}

func (s *Service) checkReadOnly(ctx context.Context) error {
	if readOnly, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok && readOnly {
		return fmt.Errorf("shop is in read-only mode")
	}
	return nil
}

func (s *Service) GetStoreDetails(ctx context.Context) (*insales.StoreInfo, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetStoreInfo(ctx)
}

// === ТОВАРЫ ===

func (s *Service) GetProducts(ctx context.Context, page, perPage int) ([]insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListProducts(ctx, page, perPage)
}

func (s *Service) GetProduct(ctx context.Context, productID string) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetProduct(ctx, productID)
}

func (s *Service) CreateProduct(ctx context.Context, title, price string, args map[string]any) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateProduct(ctx, title, price, args)
}

func (s *Service) UpdateProduct(ctx context.Context, productID string, args map[string]any) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateProduct(ctx, productID, args)
}

func (s *Service) DeleteProduct(ctx context.Context, productID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteProduct(ctx, productID)
}

// === ЗАКАЗЫ ===

func (s *Service) GetOrders(ctx context.Context, page, perPage int) ([]insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListOrders(ctx, page, perPage)
}

func (s *Service) GetOrderDetails(ctx context.Context, orderID string) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetOrder(ctx, orderID)
}

func (s *Service) UpdateOrder(ctx context.Context, orderID string, args map[string]any) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateOrder(ctx, orderID, args)
}

func (s *Service) CreateOrder(ctx context.Context, clientID int, args map[string]any) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateOrder(ctx, clientID, args)
}

// === КЛИЕНТЫ ===

func (s *Service) GetClients(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListClients(ctx, page, perPage)
}

func (s *Service) GetClient(ctx context.Context, clientID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetClient(ctx, clientID)
}

func (s *Service) CreateClient(ctx context.Context, email string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateClient(ctx, email, args)
}

func (s *Service) UpdateClient(ctx context.Context, clientID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateClient(ctx, clientID, args)
}

func (s *Service) DeleteClient(ctx context.Context, clientID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteClient(ctx, clientID)
}

// === КАТЕГОРИИ ===

func (s *Service) GetCategories(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCategories(ctx, page, perPage)
}

func (s *Service) GetCategory(ctx context.Context, categoryID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCategory(ctx, categoryID)
}

func (s *Service) CreateCategory(ctx context.Context, title string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCategory(ctx, title, args)
}

func (s *Service) UpdateCategory(ctx context.Context, categoryID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCategory(ctx, categoryID, args)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCategory(ctx, categoryID)
}

// === КОЛЛЕКЦИИ ===

func (s *Service) GetCollections(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCollections(ctx, page, perPage)
}

func (s *Service) GetCollection(ctx context.Context, collectionID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCollection(ctx, collectionID)
}

func (s *Service) CreateCollection(ctx context.Context, title string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCollection(ctx, title)
}

func (s *Service) UpdateCollection(ctx context.Context, collectionID, title string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCollection(ctx, collectionID, title)
}

func (s *Service) DeleteCollection(ctx context.Context, collectionID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCollection(ctx, collectionID)
}

// === ИЗОБРАЖЕНИЯ ===

func (s *Service) ListImages(ctx context.Context, productID string) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListImages(ctx, productID)
}

func (s *Service) AddImage(ctx context.Context, productID, url string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.AddImage(ctx, productID, url, args)
}

func (s *Service) DeleteImage(ctx context.Context, productID, imageID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteImage(ctx, productID, imageID)
}

// === МЕТАПОЛЯ ===

func (s *Service) ListMetaFields(ctx context.Context, args map[string]any) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListMetaFields(ctx, args)
}

func (s *Service) CreateMetaField(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateMetaField(ctx, args)
}

func (s *Service) UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateMetaField(ctx, metaFieldID, args)
}

func (s *Service) DeleteMetaField(ctx context.Context, metaFieldID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteMetaField(ctx, metaFieldID)
}

// === ВЕБХУКИ ===

func (s *Service) ListWebhooks(ctx context.Context) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListWebhooks(ctx)
}

func (s *Service) CreateWebhook(ctx context.Context, event, url string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateWebhook(ctx, event, url, args)
}

func (s *Service) UpdateWebhook(ctx context.Context, webhookID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateWebhook(ctx, webhookID, args)
}

func (s *Service) DeleteWebhook(ctx context.Context, webhookID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteWebhook(ctx, webhookID)
}

// === СИСТЕМНЫЕ ===

func (s *Service) ListCurrencies(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCurrencies(ctx)
}

func (s *Service) ListDeliveryVariants(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListDeliveryVariants(ctx)
}

func (s *Service) ListPaymentGateways(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPaymentGateways(ctx)
}

// === СТРАНИЦЫ ===

func (s *Service) ListPages(ctx context.Context, page, perPage int) ([]insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPages(ctx, page, perPage)
}

func (s *Service) GetPage(ctx context.Context, pageID string) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetPage(ctx, pageID)
}

func (s *Service) CreatePage(ctx context.Context, title, content, handle string) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreatePage(ctx, title, content, handle)
}

func (s *Service) UpdatePage(ctx context.Context, pageID string, args map[string]any) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdatePage(ctx, pageID, args)
}

func (s *Service) DeletePage(ctx context.Context, pageID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeletePage(ctx, pageID)
}

// === АКСЕССУАРЫ ===

func (s *Service) ListAccessories(ctx context.Context, page, perPage int) ([]insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListAccessories(ctx, page, perPage)
}

func (s *Service) GetAccessory(ctx context.Context, accessoryID string) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAccessory(ctx, accessoryID)
}

func (s *Service) CreateAccessory(ctx context.Context, title string) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateAccessory(ctx, title)
}

func (s *Service) UpdateAccessory(ctx context.Context, accessoryID string, args map[string]any) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateAccessory(ctx, accessoryID, args)
}

func (s *Service) DeleteAccessory(ctx context.Context, accessoryID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteAccessory(ctx, accessoryID)
}

// === ВАРИАНТЫ ===

func (s *Service) ListVariants(ctx context.Context, productID string) ([]insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListVariants(ctx, productID)
}

func (s *Service) GetVariant(ctx context.Context, productID, variantID string) (*insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetVariant(ctx, productID, variantID)
}

func (s *Service) UpdateVariant(ctx context.Context, productID, variantID string, args map[string]any) (*insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateVariant(ctx, productID, variantID, args)
}

func (s *Service) DeleteVariant(ctx context.Context, productID, variantID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteVariant(ctx, productID, variantID)
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func (s *Service) ListPropertyValues(ctx context.Context, accessoryID string) ([]insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPropertyValues(ctx, accessoryID)
}

func (s *Service) CreatePropertyValue(ctx context.Context, accessoryID, title string) (*insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreatePropertyValue(ctx, accessoryID, title)
}

func (s *Service) UpdatePropertyValue(ctx context.Context, propertyValueID string, args map[string]any) (*insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdatePropertyValue(ctx, propertyValueID, args)
}

func (s *Service) DeletePropertyValue(ctx context.Context, propertyValueID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeletePropertyValue(ctx, propertyValueID)
}

// === МЕНЮ ===

func (s *Service) ListMenus(ctx context.Context) ([]insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListMenus(ctx)
}

func (s *Service) GetMenu(ctx context.Context, menuID string) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetMenu(ctx, menuID)
}

func (s *Service) CreateMenu(ctx context.Context, title, handle string) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateMenu(ctx, title, handle)
}

func (s *Service) UpdateMenu(ctx context.Context, menuID string, args map[string]any) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateMenu(ctx, menuID, args)
}

func (s *Service) DeleteMenu(ctx context.Context, menuID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteMenu(ctx, menuID)
}

// === РЕДИРЕКТЫ ===

func (s *Service) ListRedirects(ctx context.Context) ([]insales.Redirect, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListRedirects(ctx)
}

func (s *Service) CreateRedirect(ctx context.Context, url, target string) (*insales.Redirect, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateRedirect(ctx, url, target)
}

func (s *Service) DeleteRedirect(ctx context.Context, redirectID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteRedirect(ctx, redirectID)
}

// === БЛОГИ ===

func (s *Service) ListBlogs(ctx context.Context) ([]insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListBlogs(ctx)
}

func (s *Service) GetBlog(ctx context.Context, blogID string) (*insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetBlog(ctx, blogID)
}

func (s *Service) CreateBlog(ctx context.Context, title, handle string) (*insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateBlog(ctx, title, handle)
}

func (s *Service) DeleteBlog(ctx context.Context, blogID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteBlog(ctx, blogID)
}

func (s *Service) ListBlogPosts(ctx context.Context, blogID string) ([]insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListBlogPosts(ctx, blogID)
}

func (s *Service) GetBlogPost(ctx context.Context, blogID, postID string) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetBlogPost(ctx, blogID, postID)
}

func (s *Service) CreateBlogPost(ctx context.Context, blogID, title, bodyHTML string) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateBlogPost(ctx, blogID, title, bodyHTML)
}

func (s *Service) UpdateBlogPost(ctx context.Context, blogID, postID string, args map[string]any) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateBlogPost(ctx, blogID, postID, args)
}

func (s *Service) DeleteBlogPost(ctx context.Context, blogID, postID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteBlogPost(ctx, blogID, postID)
}
